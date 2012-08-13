package via

import (
	"fmt"
	"github.com/str1ngs/util/json"
	"path"
	"strings"
	"time"
)

type Plan struct {
	Name         string
	Version      string
	Url          string
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
	return path.Join(config.Plans, p.Name+".json")
}

// TODO: make this atomic
func (p *Plan) Save() (err error) {
	return json.Write(p, p.Path())
}

func ReadPlan(n string) (plan *Plan, err error) {
	plan = &Plan{Name: n}
	err = json.Read(plan, plan.Path())
	return plan, err
}

func ReadPath(p string) (plan *Plan, err error) {
	s := strings.Split(path.Base(p), ".")
	if len(s) != 2 {
		return nil, fmt.Errorf("expected {name}.json got %v", s)
	}
	return ReadPlan(s[0])
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
