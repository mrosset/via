package via

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"util/file"
	"util/json"
)

type Plan struct {
	Name    string
	Version string
	Url     string
}

func (this *Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", this.Name, this.Version)
}

func (this *Plan) Print() {
	pp := func(f, v string) {
		fmt.Printf("%-10.10s = %s\n", f, v)
	}
	pp("Name", this.Name)
	pp("Version", this.Version)
	pp("Url", this.Url)
}

func (this *Plan) File() string {
	return filepath.Join(config.Plans, this.Name+".json")
}

func (this *Plan) Save() (err error) {
	return json.Write(this, this.File())
}

func ReadPlan(name string) (plan *Plan, err error) {
	plan = &Plan{Name: name}
	if !file.Exists(path.Join(config.Plans, name+".json")) {
		return nil, fmt.Errorf("Could not find plan %s.", name)
	}
	err = json.Read(plan, plan.File())
	return plan, err
}

func (this *Plan) PackageFile() string {
	return fmt.Sprintf("%s-%s-%s.tar.gz", this.NameVersion(), runtime.GOOS, runtime.GOARCH)
}
