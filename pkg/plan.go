package via

import (
	"fmt"
	"github.com/str1ngs/util/console"
	"github.com/str1ngs/util/human"
	"github.com/str1ngs/util/json"
	"path"
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

// Provides json.Template interface
func (p *Plan) SetTemplate(i interface{}) {
	c := *i.(*Plan)
	p.template = &c
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
		console.Println(p.Name, human.ByteSize(p.Size))
	}
	console.Flush()
}

type Plan struct {
	Name         string
	Version      string
	Url          string
	Group        string
	StageDir     string
	Inherit      string
	BuildInStage bool
	Date         time.Time
	Size         int64
	SubPackages  []string
	Depends      []string
	Flags        Flags
	Patch        []string
	Build        []string
	Package      []string
	PostInstall  []string
	Remove       []string
	Files        []string

	// internal
	template *Plan
}

func (p *Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Version)
}

func (p *Plan) Path() string {
	return path.Join(config.Plans, p.Group, p.Name+".json")
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
	err = json.Execute(plan)
	if err != nil {
		return nil, err
	}
	return plan, err
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

func (p Plan) SourceFile() string {
	return join(cache.Srcs(), path.Base(p.Url))
}

func (p Plan) SourcePath() string {
	return path.Join(cache.Srcs(), path.Base(p.GetUrl()))
}

func (p Plan) GetUrl() string {
	return p.Url
}

func (p Plan) GetBuildDir() string {
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
	branch, _ := config.Branch()
	return join(config.Repo, branch, p.PackageFile())
}

func (p Plan) stageDir() string {
	if p.StageDir != "" {
		return p.StageDir
	}
	return p.NameVersion()
}
