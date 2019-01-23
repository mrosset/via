package main

import (
	"fmt"
	"github.com/mrosset/via/pkg"
	"gopkg.in/urfave/cli.v2"
)

func init() {
	app.Commands = append(app.Commands, installCommand)
}

var installCommand = &cli.Command{
	Name:  "install",
	Usage: "install a package",
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
			Name:  "b",
			Value: false,
			Usage: "use experimental batch installer",
		},
	},
	ShellComplete: func(ctx *cli.Context) {
		plans, err := via.GetPlans()
		if err != nil {
			elog.Println(err)
			return
		}
		if ctx.NArg() > 0 {
			return
		}
		for _, p := range plans {
			fmt.Printf("%s ", p.Name)
		}
	},
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("b") {
			return batch(ctx)
		}
		if !ctx.Args().Present() {
			return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
		}

		via.Root(ctx.String("r"))

		for _, arg := range ctx.Args().Slice() {
			p, err := via.NewPlan(config, arg)
			if err != nil {
				return err
			}
			return via.Install(config, p.Name)
		}
		return nil
	},
}
