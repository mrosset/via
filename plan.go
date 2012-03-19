package via

import (
	"runtime"
	"util"
	"util/json"
)

type Plan struct {
	Name    string
	Version string
	Url     string
}

func (this *Plan) NameVersion() string {
	return util.Sprintf("%s-%s", this.Name, this.Version)
}

func (this *Plan) File() string {
	return plans.File(this.Name + ".json")
}

func (this *Plan) Save() (err error) {
	return json.Write(this, this.File())
}

func ReadPlan(name string) (plan *Plan, err error) {
	plan = &Plan{Name: name}
	err = json.Read(plan, plan.File())
	return plan, err
}

func (this *Plan) PackageFile() string {
	return util.Sprintf("%s-%s-%s.tar.gz", this.NameVersion(), runtime.GOOS, runtime.GOARCH)
}
