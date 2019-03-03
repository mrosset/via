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
		{
			Name:   "env",
			Usage:  "prints out via build env in shell format",
			Action: env,
		},
		{
			Name:          "stage",
			Usage:         "downloads and stages Plans source files",
			Action:        stage,
			ShellComplete: planArgCompletion,
		},
		{
			Name:   "repo",
			Usage:  "recreates file db",
			Action: repo,
		},
		{
			Name:   "diff",
			Usage:  "diff's plan working directory against git HEAD",
			Action: diff,
		},
		{
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
			Action: strap,
		},
		{
			Name:   "daemon",
			Usage:  "starts build daemon",
			Action: daemon,
		},
		{
			Name:   "hash",
			Usage:  "DEV ONLY sync the plans Oid with binary banch",
			Action: notimplemented,
		},
		{
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
		{
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
		{
			Name:   "fix",
			Usage:  "DEV ONLY used to mass modify plans",
			Action: notimplemented,
		},
		{
			Name:  "reset",
			Usage: "resets entire branch's plans",
			Description: `Resets an entire Branch's dynamic plan meta data. This Essential puts the branch in a state as if no plans were built. Its also resets any repo data.

This is useful for creating a new branch that either has another config or to bootstrap a Branch for another operating system or CPU architecture.`,
			Action: reset,
		},
		{
			Name:   "test",
			Usage:  "installs devel group into a temp directory",
			Action: test,
		},
		{
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
	var (
		sfmt  = "%-10s %-10s %-10s\n"
		files []string
		err   error
	)

	if files, err = via.PlanFiles(config); err != nil {
		return err
	}

	for _, f := range files {
		plan, err := via.ReadPath(f)
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
			if err := via.WritePlan(config, plan); err != nil {
				return err
			}
		}
	}
	return nil
}

func strap(ctx *cli.Context) error {

	dplan, err := via.NewPlan(config, "devel")

	if err != nil {
		return err
	}

	via.Debug(ctx.Bool("d"))

	for _, p := range dplan.ManualDepends {
		plan, err := via.NewPlan(config, p)
		if err != nil {
			return err
		}
		if ctx.Bool("m") {
			plan.IsRebuilt = false
			via.WritePlan(config, plan)
			continue
		}
		if plan.IsRebuilt {
			fmt.Printf(lfmt, "rebuilt", plan.NameVersion())
			continue
		}
		if err := insDeps(plan); err != nil {
			return err
		}
		if err := via.Clean(config, plan); err != nil {
			return err
		}

		b := via.NewBuilder(config, plan)
		if err := b.BuildSteps(); err != nil {
			return err
		}
	}
	return nil

}

func insDeps(plan *via.Plan) error {
	for _, p := range plan.BuildDepends {
		plan, err := via.NewPlan(config, p)
		if err != nil {
			return err
		}
		b := via.NewBatch(config, os.Stdout)
		b.Walk(plan)
		errs := b.Install()
		if len(errs) > 0 {
			return errs[0]
		}
	}
	return nil
}

func env(ctx *cli.Context) error {
	for _, k := range []string{"CFLAGS", "LDFLAGS", "CXXFLAGS"} {
		fmt.Printf("export \"%s\"\n", config.Expand().Env.Value(k))
	}
	return nil
}

func reset(ctx *cli.Context) error {
	var (
		files []string
		err   error
	)
	if files, err = via.PlanFiles(config); err != nil {
		return err
	}
	for _, path := range files {
		plan, err := via.ReadPath(path)
		if err != nil {
			return err
		}
		if plan.Group == "builtin" || plan.Group == "groups" {
			continue
		}
		plan.Cid = ""
		plan.IsRebuilt = false
		plan.Date = time.Now()
		plan.BuildTime = 0
		plan.Size = 0
		if err := via.WritePlan(config, plan); err != nil {
			return err
		}
	}
	return via.RepoCreate(config)
}

func stage(ctx *cli.Context) error {
	arg := ctx.Args().First()
	plan, err := via.NewPlan(config, arg)
	if err != nil {
		return err
	}
	b := via.NewBuilder(config, plan)
	return b.Stage()
}

func test(ctx *cli.Context) error {
	var (
		batch = via.NewBatch(config, os.Stdout)
		plan  = &via.Plan{}
		root  = ""
		err   error
	)
	if root, err = ioutil.TempDir("", "via-test"); err != nil {
		return err
	}
	defer os.RemoveAll(root)
	config.Root = via.Path(root)
	config.Repo = via.Path(root).Join("repo").ToRepo()
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
}
