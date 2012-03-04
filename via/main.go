package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
		checkf(err)
		checkf(via.DownloadSrc(plan))
		checkf(via.Stage(plan))
		checkf(via.Build(plan))
		checkf(via.MakeInstall(plan))
		checkf(via.Package(plan))
		fmt.Printf("%-20s %s\n", plan.NameVersion(), time.Now().Sub(start))
	}
}

func sign(args []string) {
	tarballs, err := filepath.Glob("*.tar.gz")
	checkf(err)
	checkf(via.Sign(tarballs))
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
