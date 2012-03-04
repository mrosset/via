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
)

func main() {
	flag.Parse()
	via.Verbose = *verbose
	util.Verbose = *verbose
	via.InitConfig()
	cmd := flag.Arg(0)
	args := flag.Args()[1:]
	switch cmd {
	case "build":
		build(args)
	case "install":
		install(args)
	case "remove":
		remove(args)
	case "sign":
		sign(args)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func build(args []string) {
	for _, arg := range args {
		start := time.Now()
		plan, err := via.ReadPlan(arg)
		util.CheckFatal(err)
		util.CheckFatal(via.DownloadSrc(plan))
		util.CheckFatal(via.Stage(plan))
		util.CheckFatal(via.Build(plan))
		util.CheckFatal(via.MakeInstall(plan))
		util.CheckFatal(via.Package(plan))
		util.CheckFatal(via.Sign(plan.PackageFile()))
		fmt.Printf("%-20s %s\n", plan.NameVersion(), time.Now().Sub(start))
	}
}

func sign(args []string) {
	for _, arg := range args {
		err := via.Sign(args)
		util.CheckFatal(err)
	}
}

func install(args []string) {
	for _, arg := range args {
		util.CheckFatal(via.Install(arg))
	}
}

func remove(args []string) {
	for _, arg := range args {
		util.CheckFatal(via.Remove(arg))
	}
}
