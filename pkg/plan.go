package via

import (
	"fmt"
	"github.com/str1ngs/util/console"
	"github.com/str1ngs/util/human"
	"github.com/str1ngs/util/json"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"
)

type PlanSlice []*Plan

// Returns a PlanSlice of all Plans in config.Plans
func NewPlanSlice() (PlanSlice, error) {
	pf, err := PlanFiles()
	if err != nil {
		return nil, err
	}
	plans := PlanSlice{}
	for _, f := range pf {
		p, _ := ReadPath(f)
		plans = append(plans, p)
	}
	return plans, nil
}

// Returns a copy of this PlanSlice sorted by
// field Size.
func (ps PlanSlice) SortSize() PlanSlice {
	nps := append(PlanSlice{}, ps...)
	sort.Sort(Size(nps))
	return nps
}

// Prints this slice to console.
// TODO: use template
func (ps PlanSlice) Print() {
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
	e, err := filepath.Glob(join(config.Plans, "*", n+".json"))
	if err != nil {
		return "", err
	}
	if len(e) != 1 {
		return "", fmt.Errorf("expected 1 plan found %d.", len(e))
	}
	return e[0], nil
}

func FindPlan(n string) (plan *Plan, err error) {
	p, err := FindPlanPath(n)
	if err != nil {
		return nil, err
	}
	return ReadPath(p)
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

func (p *Plan) Expand(s string) string {
	fn := func(s string) string {
		switch s {
		case "N":
			return p.Name
		case "V":
			return p.Version
		case "GM":
			return "http://mirrors.kernel.org/gnu"
		}
		return ""
	}
	switch s {
	case "Url":
		s = p.Url
	}
	return os.Expand(s, fn)
}

func (p Plan) stageDir() string {
	if p.StageDir != "" {
		return p.StageDir
	}
	return p.NameVersion()
}
