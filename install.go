package main

import (
	"fmt"
	"github.com/mrosset/util/file"
	"github.com/mrosset/via/pkg"
	"gopkg.in/urfave/cli.v2"
	"os"
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
					Value: config.Root,
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
			Usage:   "upgrade currently installed packages with newer build or versions",
			Action: func(ctx *cli.Context) error {
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
			},
		},
	}
)

// FIXME: this function is deprecated and should be replaced with batch
func install(ctx *cli.Context) error {
	if ctx.Bool("b") {
		return batch(ctx)
	}
	if !ctx.Args().Present() {
		return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
	}

	config.Root = ctx.String("r")
	if !file.Exists(ctx.String("r")) {
		if err := os.MkdirAll(ctx.String("r"), 0755); err != nil {
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

func planArgCompletion(ctx *cli.Context) {
	plans, err := via.GetPlans(config)
	if err != nil {
		elog.Println(err)
		return
	}
	for _, p := range plans {
		fmt.Printf("%s ", p.Name)
	}
}
