package via

import (
	"fmt"
	"github.com/str1ngs/util/json"
	"path"
	"path/filepath"
	"time"
)

type Plan struct {
	Name         string
	Version      string
	Url          string
	Group        string
	StageDir     string
	BuildInStage bool
	Date         time.Time
	Size         int64
	Depends      []string
	Flags        Flags
	Build        []string
	Package      []string
	PostInstall  []string
	Remove       []string
	Files        []string
}

func (p *Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Version)
}

func (p *Plan) Print() {
	pp := func(f, v string) {
		fmt.Printf("%-10.10s = %s\n", f, v)
	}
	pp("Name", p.Name)
	pp("Version", p.Version)
	pp("Url", p.Url)
	pp("Flags", p.Flags.String())
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

func (p Plan) stageDir() string {
	if p.StageDir != "" {
		return p.StageDir
	}
	return p.NameVersion()
}
