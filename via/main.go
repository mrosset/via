package main

import (
	"code.google.com/p/via"
	"flag"
	"fmt"
	"github.com/str1ngs/util"
	"github.com/str1ngs/util/console/command"
	"log"
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
	command.Add("remove", remove, "remove package")
	command.Add("create", create, "create plan from URL")
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func create() error {
	for _, arg := range command.Args() {
		err := via.Create(arg)
		if err != nil {
			return err
		}
	}
	return nil
}

func build() error {
	for _, arg := range command.Args() {
		start := time.Now()
		plan, err := via.ReadPlan(arg)
		defer fmt.Printf("%-20s %s\n", plan.NameVersion(), time.Now().Sub(start))
		checkf(err)
		checkf(via.BuildSteps(plan))
	}
	return nil
}

func install() error {
	return command.ArgsDo(via.Install)
}

func remove() error {
	return command.ArgsDo(via.Remove)
}
