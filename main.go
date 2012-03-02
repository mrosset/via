package main

import (
	"flag"
	"os"
	"util"
)

var (
	verbose = flag.Bool("-v", true, "verbose output")
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "build":
		InitConfig()
		build(flag.Args()[1:])
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func build(args []string) {
	for _, arg := range args {
		plan, err := ReadPlan(arg)
		util.CheckFatal(err)
		util.CheckFatal(DownloadSrc(plan))
		util.CheckFatal(Stage(plan))
		util.CheckFatal(Build(plan))
		util.CheckFatal(Install(plan))
	}
}
