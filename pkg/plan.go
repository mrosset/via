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

type Plans []*Plan

// Returns a PlanSlice of all Plans in config.Plans
func GetPlans() (Plans, error) {
	pf, err := PlanFiles()
	if err != nil {
		return nil, err
	}
	plans := Plans{}
	for _, f := range pf {
		p, _ := ReadPath(f)
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
	SourceCid     string
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
}

func (p *Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Version)
}

func (p *Plan) Path() string {
	return filepath.Join(config.Plans, p.Group, p.Name+".json")
}

// TODO: make this atomic
func (p *Plan) Save() (err error) {
	return json.Write(p, p.Path())
}

func FindPlanPath(n string) (string, error) {
	glob := join(config.Plans, "*", n+".json")
	e, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}
	if len(e) != 1 {
		return "", fmt.Errorf("expected 1 plan found %d.", len(e))
	}
	return e[0], nil
}

func NewPlan(n string) (plan *Plan, err error) {
	path, err := FindPlanPath(n)
	if err != nil {
		return nil, err
	}
	plan, err = ReadPath(path)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func ReadPath(p string) (plan *Plan, err error) {
	plan = new(Plan)
	err = json.Read(plan, p)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (p *Plan) PackageFile() string {
	return fmt.Sprintf("%s-%s-%s.tar.gz", p.NameVersion(), config.OS, config.Arch)
}

func (p *Plan) SourceFile() string {
	return join(filepath.Base(p.Expand().Url))
}

func (p *Plan) SourcePath() string {
	return filepath.Join(cache.Sources(), filepath.Base(p.Expand().Url))
}

func (p Plan) BuildDir() string {
	if p.BuildInStage {
		return join(cache.Stages(), p.stageDir())
	}
	return join(cache.Builds(), p.NameVersion())
}

func (p Plan) GetStageDir() string {
	if p.SourceCid != "" {
		return join(cache.Stages(), p.SourceCid)
	}
	return join(cache.Stages(), p.stageDir())
}

func (p Plan) PackagePath() string {
	return join(config.Repo, "repo", p.PackageFile())
}

func (p Plan) stageDir() string {
	if p.StageDir != "" {
		return p.StageDir
	}
	return p.NameVersion()
}
