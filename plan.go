package via

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type Plan struct {
	Name    string "name"
	Version string "version"
	Mirror  string "mirror"
	File    string "file"
}

func (this *Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", this.Name, this.Version)
}

func (this *Plan) Url() string {
	return fmt.Sprintf("%s/%s", this.Mirror, this.File)
}

func (this *Plan) Print() {
	pp := func(f, v string) {
		fmt.Printf("%-10.10s = %s\n", f, v)
	}
	pp("Name", this.Name)
	pp("Version", this.Version)
	pp("File", this.File)
	pp("Mirror", this.Mirror)
}

func (this *Plan) Save() (err error) {
	return WriteJson(this, filepath.Join(config.Plans, this.Name+".json"))
}

func ReadPlan(name string) (plan *Plan, err error) {
	plan = &Plan{}
	plan, err = ReadJson(name)
	return plan, err
}

func (this *Plan) PackageFile() string {
	return fmt.Sprintf("%s-%s-%s.tar.gz", this.NameVersion(), runtime.GOOS, runtime.GOARCH)
}
