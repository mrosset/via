package main

import (
	"code.google.com/p/via"
	"flag"
	"fmt"
	"github.com/str1ngs/util"
	"github.com/str1ngs/util/console/command"
	"github.com/str1ngs/util/json"
	"log"
	"os"
	"os/exec"
)

var (
	verbose  = flag.Bool("v", false, "verbose output")
	finstall = flag.Bool("i", false, "install package after build")
)

func main() {
	via.Verbose = *verbose
	util.Verbose = *verbose
	command.Add("build", build, "build plan")
	command.Add("clean", clean, "clean build dir")
	command.Add("create", create, "create plan from URL")
	command.Add("edit", edit, "calls EDITOR to edit plan")
	command.Add("files", files, "lists files")
	command.Add("install", install, "install package")
	command.Add("lint", lint, "lint plans")
	command.Add("list", list, "list all plans")
	command.Add("pack", pack, "package plan")
	command.Add("remove", remove, "remove package")
	command.Add("show", xshow, "displays plan to stdout")
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func edit() error {
	editor := os.Getenv("EDITOR")
	plan, err := via.ReadPlan(command.Args()[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(editor, plan.Path())
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
		plan, err := via.ReadPlan(arg)
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
		plan, err := via.ReadPlan(arg)
		if err != nil {
			return err
		}
		err = via.Package(plan)
		if err != nil {
			return err
		}
	}
	return nil
}

func files() error {
	for _, arg := range command.Args() {
		plan, err := via.ReadPlan(arg)
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
		plan, err := via.ReadPlan(arg)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Clean(&plan, os.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}

func clean() error {
	return command.ArgsDo(via.Clean)
}

func list() error {
	via.List()
	return nil
}
