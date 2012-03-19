package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"util"
	"via"
)

var (
	verbose = flag.Bool("v", false, "verbose output")
	checkf  = util.CheckFatal
)

func main() {
	flag.Parse()
	via.Verbose = *verbose
	util.Verbose = *verbose
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	cmd := flag.Arg(0)
	args := flag.Args()[1:]
	switch cmd {
	case "build":
		build(args)
	case "install":
		install(args)
	case "remove":
		remove(args)
	case "create":
		create(args)
	case "list":
		list(args)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func create(args []string) {
	for _, arg := range args {
		checkf(via.Create(arg))
	}
}

func build(args []string) {
	for _, arg := range args {
		start := time.Now()
		plan, err := via.ReadPlan(arg)
		defer fmt.Printf("%-20s %s\n", plan.NameVersion(), time.Now().Sub(start))
		checkf(err)
		checkf(via.BuildSteps(plan))
	}
}

func install(args []string) {
	for _, arg := range args {
		checkf(via.Install(arg))
	}
}

func remove(args []string) {
	for _, arg := range args {
		checkf(via.Remove(arg))
	}
}

func list(args []string) {
	for _, arg := range args {
		checkf(via.List(arg))
	}
}
