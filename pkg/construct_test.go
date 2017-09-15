package via

import (
	"fmt"
)

var (
	tConstruct *Construct
)

func init() {
	plan, _ := NewPlan("bash")
	tConstruct = NewConstruct(config, plan)
	tConstruct.Cache = "testdata/cache"
}
func ExamplePaths() {
	fmt.Println(tConstruct.PackageFilePath())
	fmt.Println(tConstruct.PlanSourcePath())
	fmt.Println(tConstruct.BuildPath())
	fmt.Println(tConstruct.PlanStagePath())
	fmt.Println(tConstruct.Cache.Stages())
	fmt.Println(tConstruct.Cache.Sources())
	fmt.Println(tConstruct.Cache.Builds())
	fmt.Println(tConstruct.Cache.Packages())
	// output:
	// /home/mrosset/gocode/src/github.com/mrosset/via/publish/repo/bash-4.3-linux-x86_64.tar.gz
	// testdata/cache/src/bash-4.3.tar.gz
	// testdata/cache/bld/bash-4.3
	// testdata/cache/stg/bash-4.3
	// testdata/cache/stg
	// testdata/cache/src
	// testdata/cache/bld
	// testdata/cache/pkg
}
