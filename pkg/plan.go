package via

import (
	"fmt"
	"github.com/mrosset/util/console"
	"github.com/mrosset/util/human"
	"github.com/mrosset/util/json"
	"path/filepath"
	"sort"
	"time"
)

// Provides a slice of plans
type Plans []*Plan

// Returns a Plan slice of all Plans in config.Plans
func GetPlans(config *Config) (Plans, error) {
	pf, err := PlanFiles(config)
	if err != nil {
		return nil, err
	}
	plans := Plans{}
	for _, f := range pf {
		p, _ := ReadPath(config, f)
		plans = append(plans, p)
	}
	return plans, nil
}

// Returns a copy of this PlanSlice sorted by
// field Size.
func (ps Plans) SortSize() Plans {
	nps := append(Plans{}, ps...)
	sort.Sort(Size(nps))
	return nps
}

// Prints this slice to console.
// TODO: use template
func (ps Plans) Print() {
	for _, p := range ps {
		console.Println(p.NameVersion(), human.ByteSize(p.Size))
	}
	console.Flush()
}

// Returns a slice of plan names
func (ps Plans) Slice() []string {
	s := []string{}
	for _, p := range ps {
		s = append(s, p.Name)
	}
	return s
}

func (ps Plans) Contains(plan *Plan) bool {
	for _, p := range ps {
		if p.Name == plan.Name {
			return true
		}
	}
	return false
}

// Returns a expanded Plan template.
func (p *Plan) Expand() *Plan {
	o := new(Plan)
	err := json.Parse(o, p)
	if err != nil {
		panic(err)
	}
	return o
}

type Plan struct {
	Name          string
	Version       string
	Url           string
	Group         string
	StageDir      string
	Inherit       string
	Cid           string
	BuildInStage  bool
	IsRebuilt     bool
	BuildTime     time.Duration
	Date          time.Time
	Size          int64
	SubPackages   []string
	AutoDepends   []string
	ManualDepends []string
	BuildDepends  []string
	Flags         Flags
	Patch         []string
	Build         []string
	Package       []string
	PostInstall   []string
	Remove        []string
	Files         []string
	config        *Config
}

func (p *Plan) Depends() []string {
	return append(p.AutoDepends, p.ManualDepends...)
}

func (p *Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Version)
}

func FindPlanPath(config *Config, name string) (string, error) {
	glob := join(config.Plans, "*", name+".json")
	e, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}
	if len(e) != 1 {
		return "", fmt.Errorf("%s: expected 1 plan found %d.", name, len(e))
	}
	return e[0], nil
}

func NewPlan(config *Config, name string) (plan *Plan, err error) {
	path, err := FindPlanPath(config, name)
	if err != nil {
		return nil, err
	}
	plan, err = ReadPath(config, path)
	if err != nil {
		return nil, err
	}
	plan.config = config
	return plan, nil
}

func ReadPath(config *Config, path string) (plan *Plan, err error) {
	plan = new(Plan)
	err = json.Read(plan, path)
	if err != nil {
		return nil, err
	}
	plan.config = config
	return plan, nil
}

func (p *Plan) PackageFile() string {
	if p.Cid == "" {
		return fmt.Sprintf("%s-%s-%s.tar.gz", p.NameVersion(), p.config.OS, p.config.Arch)
	}
	return fmt.Sprintf("%s.tar.gz", p.Cid)
}

func (p *Plan) SourceFile() string {
	return join(filepath.Base(p.Expand().Url))
}

func (p Plan) PackagePath() string {
	return join(p.config.Repo, p.PackageFile())
}

func (p Plan) stageDir() string {
	if p.StageDir != "" {
		return p.StageDir
	}
	return p.NameVersion()
}
