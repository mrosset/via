package main

import (
	"code.google.com/p/via"
	"flag"
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
	command.Add("create", create, "create plan from URL")
	command.Add("edit", edit, "calls EDITOR to edit plan")
	command.Add("install", install, "install package")
	command.Add("remove", remove, "remove package")
	command.Add("lint", lint, "lint plans")
	command.Add("cat", cat, "displays plan to stdout")
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

func install() error {
	return command.ArgsDo(via.Install)
}

func remove() error {
	return command.ArgsDo(via.Remove)
}

func lint() error {
	return via.Lint()
}

func cat() error {
	for _, arg := range command.Args() {
		plan, err := via.ReadPlan(arg)
		if err != nil {
			log.Fatal(err)
		}
		err = json.WritePretty(&plan, os.Stdout)
		if err != nil {
			return err
		}
	}
	return nil
}
