package main

import (
	"fmt"
	"github.com/mrosset/util/file"
	"github.com/mrosset/via/pkg"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

func init() {
	app.Commands = append(app.Commands, installCommands...)
}

var (
	installCommands = []*cli.Command{
		&cli.Command{
			Name:    "install",
			Usage:   "install a package",
			Aliases: []string{"i"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "r",
					Value: config.Root.String(),
					Usage: "use `\"DIR\"` as root",
				},
				&cli.BoolFlag{
					Name:  "y",
					Value: true,
					Usage: "Don't prompt to install",
				},
				&cli.BoolFlag{
					Name:  "s",
					Value: false,
					Usage: "use single threaded installer",
				},
			},
			ShellComplete: planArgCompletion,
			Action:        batch,
		},
		&cli.Command{
			Name:    "upgrade",
			Aliases: []string{"u", "up"},
			Usage:   "upgrade packages with newer build or versions",
			Action:  upgrade,
		},
	}
)

func batch(ctx *cli.Context) error {
	var errors []error
	if ctx.Bool("s") {
		return install(ctx)
	}
	if !ctx.Args().Present() {
		return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
	}

	config.Root = via.Path(ctx.String("r"))

	batch := via.NewBatch(config, os.Stdout)
	for _, a := range ctx.Args().Slice() {
		p, err := via.NewPlan(config, a)
		if err != nil {
			return err
		}
		if p.Cid == "" {
			return fmt.Errorf("plan '%s' does not have a Cid. Has the plan been built?", p.Name)
		}
		if err := batch.Walk(p); err != nil {
			return err
		}
	}
	switch ctx.Bool("y") {
	case false:
		errors = batch.PromptInstall()
	case true:
		errors = batch.Install()
	}
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// FIXME: this function is deprecated and should be replaced with batch
func install(ctx *cli.Context) error {
	config.Root = via.Path(ctx.String("r"))
	if ctx.Bool("b") {
		return batch(ctx)
	}
	if !ctx.Args().Present() {
		return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
	}
	if !file.Exists(ctx.String("r")) {
		if err := os.MkdirAll(ctx.String("r"), 0755); err != nil {
			elog.Printf("could not create '%s'. %s", config.Root, err)
			return err
		}
	}
	for _, arg := range ctx.Args().Slice() {
		p, err := via.NewPlan(config, arg)
		if err != nil {
			return err
		}
		if err := via.NewInstaller(config, p).Install(); err != nil {
			return err
		}
	}
	return nil
}

func upgrade(ctx *cli.Context) error {
	up := via.NewUpgrader(config)
	p, err := up.Check()
	if err != nil {
		return err
	}
	if len(p) > 0 {
		fmt.Println("upgrading", p)
	}
	errs := up.Upgrade()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func planArgCompletion(_ *cli.Context) {
	files, err := via.PlanFiles(config)
	if err != nil {
		elog.Println(err)
		return
	}
	for _, f := range files {
		name := filepath.Base(f)
		fmt.Printf("%s ", name[:len(name)-5])
	}
}
