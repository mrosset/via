package main

import (
	"code.google.com/p/via"
	"flag"
	"fmt"
	"github.com/str1ngs/util"
	"github.com/str1ngs/util/console/command"
	"time"
)

var (
	verbose = flag.Bool("v", false, "verbose output")
	checkf  = util.CheckFatal
)

func main() {
	via.Verbose = *verbose
	util.Verbose = *verbose
	command.Add("build", build, "build plan")
	command.Add("install", install, "install package")
	command.Add("remove", install, "remove package")
	command.Add("create", install, "create plan from URL")
	command.Run()
}

func create() {
	for _, arg := range command.Args() {
		checkf(via.Create(arg))
	}
}

func build() {
	for _, arg := range command.Args() {
		start := time.Now()
		plan, err := via.ReadPlan(arg)
		defer fmt.Printf("%-20s %s\n", plan.NameVersion(), time.Now().Sub(start))
		checkf(err)
		checkf(via.BuildSteps(plan))
	}
}

func install() {
	command.ArgsDo(via.Install)
}

func remove() {
	command.ArgsDo(via.Remove)
}
