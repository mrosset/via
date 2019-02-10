package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/mrosset/via/pkg"
	"github.com/mrosset/via/pkg/upstream"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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
			Name:  "strap",
			Usage: "rebuilds each package in the devel group",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "m",
					Value: false,
					Usage: "marks package in development group for rebuild",
				},
				&cli.BoolFlag{
					Name:  "d",
					Value: false,
					Usage: "debug output",
				},
			},
			Action: func(ctx *cli.Context) error {
				return fmt.Errorf("strap command is not implemented")
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
			Action: notimplemented,
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
			Description: `Resets an entire Branch's dynamic plan meta data. This Essential puts the branch in a state as if no plans were built. Its also resets any repo data.

This is useful for creating a new branch that either has another config or to bootstrap a Branch for another operating system or CPU architecture.`,
			Action: func(ctx *cli.Context) error {
				var (
					files []string
					err   error
				)
				if files, err = via.PlanFiles(); err != nil {
					return err
				}
				for _, path := range files {
					plan, err := via.ReadPath(config, path)
					if err != nil {
						return err
					}
					plan.Cid = ""
					plan.IsRebuilt = false
					plan.Date = time.Now()
					plan.BuildTime = 0
					plan.Files = nil
					plan.Size = 0
					plan.AutoDepends = nil
					ctx := via.NewPlanContext(config, plan)
					if err := ctx.WritePlan(); err != nil {
						return err
					}
				}
				if err = via.RepoCreate(config); err != nil {
					return err
				}

				return nil

			},
		},
		&cli.Command{
			Name:  "test",
			Usage: "installs devel group into a temp directory",
			Action: func(ctx *cli.Context) error {
				var (
					batch = via.NewBatch(config)
					plan  = &via.Plan{}
					root  = ""
					err   error
				)
				if root, err = ioutil.TempDir("", "via"); err != nil {
					return err
				}
				defer os.RemoveAll(root)
				config.Root = root
				config.Repo = filepath.Join(root, "repo")
				if plan, err = via.NewPlan(config, "devel"); err != nil {
					return err
				}
				if err = batch.Walk(plan); err != nil {
					return err
				}
				errors := batch.Install()
				if len(errors) != 0 {
					return errors[0]
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "upstream",
			Usage: "manage upstream versions and urls",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "w",
					Value: false,
					Usage: "writes new upstream versions",
				},
			},
			Action: cupstream,
		},
	},
}

func notimplemented(ctx *cli.Context) error {
	return fmt.Errorf("'%s' command is not implemented", ctx.Command.FullName())
}

func cupstream(ctx *cli.Context) error {
	files, err := via.PlanFiles()
	if err != nil {
		return err
	}
	sfmt := "%-10s %-10s %-10s\n"

	for _, f := range files {
		plan, err := via.ReadPath(config, f)
		if err != nil {
			return err
		}
		if plan.Version == "" || plan.Url == "" || plan.Cid == "" {
			continue
		}
		dir, err := url.Parse("./")
		if err != nil {
			return err
		}
		uri, err := url.Parse(plan.Expand().Url)
		if err != nil {
			return err
		}
		uri = uri.ResolveReference(dir)
		current, err := semver.ParseTolerant(plan.Version)
		if err != nil {
			fmt.Printf(sfmt, plan.Name, "error", err)
			continue
		}
		upstream, err := upstream.GnuUpstreamLatest(plan.Name, uri.String(), current)
		if err != nil {
			if oerr, ok := err.(net.Error); ok {
				fmt.Printf(sfmt, plan.Name, "error", oerr)
				continue
			}
			fmt.Printf(sfmt, plan.Name, "error", err)
			continue
		}
		if upstream != "0.0.0" {
			fmt.Printf(sfmt, plan.Name, plan.Version, upstream)
		}

		if upstream != "0.0.0" && ctx.Bool("w") {
			plan.Url = strings.Replace(plan.Url, plan.Version, "{{.Version}}", -1)
			plan.Version = upstream
			plan.IsRebuilt = false
			plan.Cid = ""
			ctx := via.NewPlanContext(config, plan)
			if err := ctx.WritePlan(); err != nil {
				return err
			}
		}
	}
	return nil
}
