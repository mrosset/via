package main

import (
	"gopkg.in/urfave/cli.v2"
)

func init() {
	app.Commands = append(app.Commands, develCommand)
}

var develCommand = &cli.Command{
	Name:    "devel",
	Usage:   "experimental and development commands",
	Aliases: []string{"dev"},
	Subcommands: []*cli.Command{
		&cli.Command{
			Name:   "repo",
			Usage:  "recreates file db",
			Action: repo,
		},
		&cli.Command{
			Name:   "diff",
			Usage:  "diff's plan working directory against git HEAD",
			Action: diff,
		},
		&cli.Command{
			Name:   "strap",
			Usage:  "rebuilds each package in the devel group",
			Action: strap,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "m",
					Value: false,
					Usage: "marks package in development group for rebuild",
				},
			},
		},
		&cli.Command{
			Name:   "daemon",
			Usage:  "starts build daemon",
			Action: daemon,
		},
		&cli.Command{
			Name:   "hash",
			Usage:  "DEV ONLY sync the plans Oid with binary banch",
			Action: hash,
		},
		&cli.Command{
			Name:  "cd",
			Usage: "prints out shell evaluate-able command to change directory. eg. eval $(via cd -s bash)",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "s",
					Value: false,
					Usage: "prints stage directory",
				},
				&cli.BoolFlag{
					Name:  "b",
					Value: false,
					Usage: "prints build directory",
				},
			},
			Action: cd,
		},
		&cli.Command{
			Name:   "plugin",
			Usage:  "execute plugin",
			Action: plugin,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "b",
					Value: false,
					Usage: "compile plugins",
				},
			},
		},
		&cli.Command{
			Name:   "edit",
			Usage:  "calls EDITOR to edit plan",
			Action: edit,
		},
		&cli.Command{
			Name:   "fix",
			Usage:  "DEV ONLY used to mass modify plans",
			Action: fix,
		},
		&cli.Command{
			Name:  "reset",
			Usage: "resets entire branch's plans",
		},
	},
}
