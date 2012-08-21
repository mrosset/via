package main

import (
	"bytes"
	"code.google.com/p/via/pkg"
	"flag"
	"fmt"
	"github.com/str1ngs/util"
	"github.com/str1ngs/util/console/command"
	"github.com/str1ngs/util/json"
	"log"
	"os"
	"os/exec"
	"path"
)

var (
	root     = flag.String("r", "/", "root directory")
	verbose  = flag.Bool("v", false, "verbose output")
	finstall = flag.Bool("i", false, "install package after build")
	fdebug   = flag.Bool("d", false, "debug output")
)

func main() {
	flag.Parse()
	via.Verbose(*verbose)
	via.Root(*root)
	util.Verbose = *verbose
	via.Debug(*fdebug)
	command.Add("build", build, "build plan")
	command.Add("clean", clean, "clean build dir")
	command.Add("create", create, "create plan from URL")
	command.Add("edit", edit, "calls EDITOR to edit plan")
	command.Add("list", list, "lists files")
	command.Add("install", install, "install package")
	command.Add("lint", lint, "lint plans")
	command.Add("repo", repo, "update repo")
	command.Add("owns", owns, "finds which package owns a file")
	command.Add("sync", sync, "fetch remote repo data")
	command.Add("search", search, "search for plans (currently lists all use grep)")
	command.Add("pack", pack, "package plan")
	command.Add("remove", remove, "remove package")
	command.Add("show", xshow, "prints plan to stdout")
	command.Add("config", config, "prints config to stdout")
	command.Add("elf", elf, "prints elf information to stdout")
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func edit() error {
	var (
		editor = os.Getenv("EDITOR")
		arg0   = command.Args()[0]
		p      = path.Join(via.GetConfig().Plans, "config.json")
		err    error
	)
	if arg0 != "config" {
		p, err = via.FindPlanPath(arg0)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(editor, p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
		plan, err := via.FindPlan(arg)
		if err != nil {
			return err
		}
		err = via.BuildSteps(plan)
		if err != nil {
			return err
		}
		if *finstall {
			err := via.Install(plan.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func pack() error {
	for _, arg := range command.Args() {
		plan, err := via.FindPlan(arg)
		if err != nil {
			return err
		}
		err = via.Package("", plan)
		if err != nil {
			return err
		}
	}
	return nil
}

func list() error {
	for _, arg := range command.Args() {
		plan, err := via.FindPlan(arg)
		if err != nil {
			return err
		}
		for _, f := range plan.Files {
			fmt.Println(f)
		}
	}
	return nil
}

func install() error {
	return command.ArgsDo(via.Install)
}

func remove() error {
	return command.ArgsDo(via.Remove)
}

func lint() error {
	return via.Lint()
}

func xshow() error {
	for _, arg := range command.Args() {
		plan, err := via.FindPlan(arg)
		if err != nil {
			log.Fatal(err)
		}
		buf := new(bytes.Buffer)
		less := exec.Command("less")
		less.Stdin = buf
		less.Stdout = os.Stdout
		less.Stderr = os.Stderr
		err = json.Clean(&plan, buf)
		if err != nil {
			fmt.Println(err)
		}
		less.Run()
	}
	return nil
}

func config() error {
	err := json.Clean(via.GetConfig(), os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func clean() error {
	return command.ArgsDo(via.Clean)
}

func elf() error {
	for _, arg := range command.Args() {
		err := via.Readelf(arg)
		if err != nil {
			return err
		}
	}
	return nil
}

func sync() error {
	return via.PlanSync()
}

func owns() error {
	rfiles, err := via.ReadRepoFiles()
	if err != nil {
		return err
	}
	for _, arg := range command.Args() {
		owner := rfiles.Owns(arg)
		if owner == "" {
			fmt.Println(arg+":", "owner not found.")
			continue
		}
		fmt.Println(owner)
	}
	return nil
}

func repo() error {
	return via.RepoCreate()
}

func search() error {
	plans, err := via.NewPlanSlice()
	if err != nil {
		return err
	}
	plans.SortSize().Print()
	return nil
}
