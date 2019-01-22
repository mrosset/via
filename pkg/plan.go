package via

import (
	"bytes"
	"fmt"
	"github.com/mrosset/util/console"
	"github.com/mrosset/util/human"
	"github.com/mrosset/util/json"
	"path/filepath"
	"sort"
	"text/template"
	"time"
)

// Provides a slice of plans
type Plans []*Plan

// Returns a Plan slice of all Plans in config.Plans
func GetPlans() (Plans, error) {
	pf, err := PlanFiles()
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

func OExpand(i interface{}, s string) string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buf, i)
	if err != nil {
		panic(err)
	}
	return buf.String()
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

func (p *Plan) Path() string {
	return filepath.Join(config.Plans, p.Group, p.Name+".json")
}

// TODO: make this atomic
func (p *Plan) Save() (err error) {
	return json.Write(p, p.Path())
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
	return fmt.Sprintf("%s-%s-%s.tar.gz", p.NameVersion(), config.OS, config.Arch)
}

func (p *Plan) SourceFile() string {
	return join(filepath.Base(p.Expand().Url))
}

func (p *Plan) SourcePath() string {
	return filepath.Join(cache.Sources(), filepath.Base(p.Expand().Url))
}

func (p Plan) BuildDir() string {
	bdir := join(cache.Builds(), p.NameVersion())
	if p.BuildInStage {
		bdir = join(cache.Stages(), p.stageDir())
	}
	return bdir
}

func (p Plan) GetStageDir() string {
	path := join(cache.Stages(), p.stageDir())
	return path
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
