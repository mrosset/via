package main

import (
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/console"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg"
	viaplugin "github.com/mrosset/via/pkg/plugin"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	goplugin "plugin"
	"sort"
	"strings"
	"text/template"
)

var (
	viapath  = via.Path(os.Getenv("GOPATH")).Join("src/github.com/mrosset/via")
	planpath = viapath.Join("plans")
	config   = readconfig()
	cfile    = viapath.Join("plans/config.json")
	viaURL   = "https://github.com/mrosset/via"
	planURL  = "https://github.com/mrosset/plans"
	viabin   = via.NewPath(os.Getenv("GOPATH")).Join("bin", "via")
	elog     = log.New(os.Stderr, "", log.Lshortfile)
	lfmt     = "%-20.20s %v\n"
	app      = &cli.App{
		Name:                  "via",
		Usage:                 "a systems package manager",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "/path/to/some",
			},
		},
	}
	commands = []*cli.Command{
		{
			Name:          "edit",
			Usage:         "calls EDITOR to edit plan",
			Action:        edit,
			ShellComplete: planArgCompletion,
		},
		{
			Name:          "build",
			Usage:         "builds a plan locally",
			Aliases:       []string{"b"},
			Action:        build,
			ShellComplete: planArgCompletion,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "c",
					Value: false,
					Usage: "clean build directory before building",
				},
				&cli.BoolFlag{
					Name:   "real",
					Value:  false,
					Hidden: true,
				},
				&cli.BoolFlag{
					Name:  "v",
					Value: true,
					Usage: "displays more information when building",
				},
				&cli.BoolFlag{
					Name:  "dd",
					Value: false,
					Usage: "build depends aswell",
				},
				&cli.BoolFlag{
					Name:  "d",
					Value: false,
					Usage: "displays debugging information when building",
				},
				&cli.BoolFlag{
					Name:  "i",
					Value: false,
					Usage: "install package after building",
				},
				&cli.BoolFlag{
					Name:  "f",
					Value: false,
					Usage: "force rebuilding",
				},
				&cli.BoolFlag{
					Name:  "u",
					Value: false,
					Usage: "force downloading of sources",
				},
				&cli.BoolFlag{
					Name:  "l",
					Value: false,
					Usage: "builds plan locally",
				},
				&cli.BoolFlag{
					Name:  "r",
					Value: false,
					Usage: "builds plan using daemon",
				},
			},
		},
		{
			Name:          "remove",
			Usage:         "uninstall package",
			Action:        remove,
			ShellComplete: planArgCompletion,
		},
		{
			Name:          "show",
			Usage:         "prints plan to stdout",
			Action:        show,
			ShellComplete: planArgCompletion,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "t",
					Value: "",
					Usage: "use go template",
				},
				&cli.BoolFlag{
					Name:  "d",
					Value: false,
					Usage: "output depends",
				},
				&cli.BoolFlag{
					Name:  "v",
					Value: false,
					Usage: "output version",
				},
			},
		},
		{
			Name:   "config",
			Usage:  "prints config to stdout",
			Action: fconfig,
		},
		{
			Name:          "list",
			Usage:         "list files for `PLAN`",
			Action:        list,
			ShellComplete: planArgCompletion,
		},
		{
			Name:   "fmt",
			Usage:  "format plans",
			Action: fmtplans,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "v",
					Value: false,
					Usage: "verbose information",
				},
			},
		},
		{
			Name:          "log",
			Usage:         "output's config.log for build",
			ShellComplete: planArgCompletion,
			Action:        plog,
		},
		{
			Name:   "elf",
			Usage:  "prints elf information to stdout",
			Action: elf,
		},
		{
			Name:   "search",
			Usage:  "lists all of the available packages",
			Action: search,
		},
		{
			Name:          "options",
			Usage:         "prints the GNU configure options for a package",
			Action:        options,
			ShellComplete: planArgCompletion,
		},
		{
			Name:   "create",
			Usage:  "creates a plan from tarball URL",
			Action: create,
		},
		{
			Name:   "pack",
			Usage:  "package plan",
			Action: pack,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "v",
					Value: false,
					Usage: "displays more information when packing",
				},
			},
		},
		{
			Name:   "debug",
			Usage:  "displays enviroment and PATH details",
			Action: debug,
		},
		{
			Name:   "owns",
			Usage:  "find which plans owns 'file'",
			Action: owns,
		},
		{
			Name:   "clean",
			Usage:  "cleans cache directory",
			Action: notimplemented,
		},
		{
			Name:          "get",
			Usage:         "download source from upstream",
			Action:        get,
			ShellComplete: planArgCompletion,
		},
		{
			Name:  "bump",
			Usage: "update version for 'PLAN'",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ver",
					Usage: "new version",
				},
			},
			ShellComplete: planArgCompletion,
			Action: func(ctx *cli.Context) error {
				if ctx.String("ver") == "" {
					return fmt.Errorf("you must specify new version with -ver")
				}
				plan, err := via.NewPlan(config, ctx.Args().First())
				if err != nil {
					return err
				}
				plan.Version = ctx.String("ver")
				return via.WritePlan(config, plan)
			},
		},
	}
)

func main() {
	app.Commands = append(app.Commands, commands...)
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	if err := app.Run(os.Args); err != nil {
		elog.Fatal(err)
	}
}

func cloneplans() error {
	//
	if planpath.Exists() {
		return nil
	}
	elog.Println("cloning plans")
	return via.CloneBranch(planpath, planURL, "x86_64-via-linux-gnu-release")
}

func readconfig() *via.Config {
	// FIXME: check this somewhere else maybe?
	if os.Getenv("GOPATH") == "" {
		elog.Fatal("GOPATH must be set")
	}
	if err := cloneplans(); err != nil {
		elog.Fatal(err)
	}
	config, err := via.NewConfig(cfile)
	if err != nil {
		elog.Fatal(err)
	}
	return config
}

func plugin(ctx *cli.Context) error {
	if ctx.Bool("b") {
		if err := viaplugin.Build(config); err != nil {
			log.Fatal(err)
		}
		return nil
	}
	if !ctx.Args().Present() {
		return fmt.Errorf("plugin requires a 'plugin' argument. see: 'via help get'")
	}
	name := ctx.Args().First()
	mod := config.Plans.Join("..", "..", "plugins", name+".so")
	plug, err := goplugin.Open(mod.String())
	if err != nil {
		elog.Fatal(err)
	}
	sym, err := plug.Lookup(strings.Title(name))
	if err != nil {
		elog.Fatal(err)
	}
	test, ok := sym.(viaplugin.Plugin)
	if !ok {
		return fmt.Errorf("symbol is not a Plugin")
	}
	test.SetConfig(config)
	return test.Execute()
}

func get(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("get requires a 'PLAN' argument. see: 'via help get'")
	}

	plan, err := via.NewPlan(config, ctx.Args().First())
	if err != nil {
		return err
	}
	return gurl.Download("./", plan.Expand().Url)
}

func daemon(_ *cli.Context) error {
	return via.StartDaemon(config)
}

func remove(ctx *cli.Context) error {
	for _, arg := range ctx.Args().Slice() {
		if err := via.Remove(config, arg); err != nil {
			return err
		}
	}
	return nil
}

func build(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("build requires a 'PLAN' argument. see: 'via help build'")
	}
	for _, arg := range ctx.Args().Slice() {
		plan, err := via.NewPlan(config, arg)
		if err != nil {
			return err
		}
		if plan.IsRebuilt && !ctx.Bool("f") {
			return fmt.Errorf("plan %s is already built", plan.Name)
		}
	}
	// if we don't have a real flag then we need to enter a contain
	if !ctx.Bool("real") && !ctx.Bool("l") && os.Getenv("INSIDE_VIA") != "true" {
		return contain(ctx)
	}
	// if r flag build package with RPC daemon
	if ctx.Bool("r") {
		return remote(ctx)
	}
	for _, arg := range ctx.Args().Slice() {
		plan, err := via.NewPlan(config, arg)
		if err != nil {
			return err
		}
		if plan.IsRebuilt && !ctx.Bool("f") {
			return fmt.Errorf("Plan is built already")
		}
		if plan.IsRebuilt && ctx.Bool("f") {
			plan.IsRebuilt = false
			via.WritePlan(config, plan)
		}
		via.Verbose(ctx.Bool("v"))
		via.Debug(ctx.Bool("d"))
		via.Update(ctx.Bool("u"))

		if ctx.Bool("c") {
			via.Clean(config, plan)
		}
		if ctx.Bool("dd") {
			if err := via.BuildDeps(config, plan); err != nil {
				return err
			}
		}
		builder := via.NewBuilder(config, plan)
		if err := builder.BuildSteps(); err != nil {
			return err
		}

		if ctx.Bool("i") {
			fmt.Printf(lfmt, "install", plan.NameVersion())
			b := via.NewBatch(config, os.Stdout)
			b.Walk(plan)
			if err := b.Install(); err != nil {
				return err[0]
			}
		}
	}
	return nil
}

func remote(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("build requires a 'PLAN' argument. see: 'via help build'")
	}
	c, err := via.Connect()
	if err != nil {
		return err
	}
	res := via.Response{}
	p, _ := via.NewPlan(config, ctx.Args().First())
	req := via.Request{*p}
	return c.Call("Builder.RpcBuild", req, &res)
}

func edit(ctx *cli.Context) error {
	var (
		editor = os.Getenv("EDITOR")
		arg0   = ctx.Args().First()
		p      = config.Plans.ConfigFile()
		err    error
	)
	if arg0 != "config" {
		p, err = via.FindPlanPath(config, arg0)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(editor, p.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	elog.Println("linting...")
	return via.FmtPlans(config)
}

func list(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("list requires a 'PLAN' argument. see: 'via help list'")
	}
	plan, err := via.NewPlan(config, ctx.Args().First())
	if err != nil {
		return err
	}
	for _, f := range plan.Files {
		fmt.Println(f)
	}
	return nil
}

func fmtplans(ctx *cli.Context) error {
	via.Verbose(ctx.Bool("v"))
	return via.FmtPlans(config)
}

func parse(input string, plan *via.Plan) error {
	var (
		t = fmt.Sprintf("%s\n", input)
	)
	tmpl, err := template.New("stdout").Parse(t)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, plan)
}

func show(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("show requires a 'PLAN' argument. see: 'via help show'")
	}
	plan, err := via.NewPlan(config, ctx.Args().First())
	if err != nil {
		elog.Fatal(err)
	}
	if ctx.String("t") != "" {
		return parse(ctx.String("t"), plan)
	}
	if ctx.Bool("d") {
		return parse("{{.AutoDepends}}", plan)
	}
	if ctx.Bool("v") {
		return parse("{{.Version}}", plan)
	}
	return json.WritePretty(&plan, os.Stdout)
}

func fconfig(_ *cli.Context) error {
	err := json.WritePretty(config, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func repo(_ *cli.Context) error {
	return via.RepoCreate(config)
}

func plog(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("show requires a 'PLAN' argument. see: 'via help log'")
	}
	b, err := via.NewBuilderByName(config, ctx.Args().First())
	if err != nil {
		return err
	}
	return file.Cat(os.Stdout,
		b.Context.BuildDir.Join("config.log").String(),
	)
}

func elf(ctx *cli.Context) error {
	fmt.Println(ctx.Args().First())
	return via.Readelf(ctx.Args().First())
}

func diff(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("diff requires a 'PLAN' argument. see: 'via help diff'")
	}
	for _, arg := range ctx.Args().Slice() {
		glob := config.Plans.Join("*", arg+".json").String()
		res, err := filepath.Glob(glob)
		if err != nil {
			return err
		}
		git := exec.Command("git", "diff", strings.Join(res, " "))
		git.Dir = config.Plans.String()
		git.Stdout = os.Stdout
		git.Stderr = os.Stderr
		err = git.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func search(_ *cli.Context) error {
	plans, err := via.GetPlans(config)
	if err != nil {
		return err
	}
	plans.SortSize().Print()
	return nil
}

func options(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("options requires a 'PLAN' argument. see: 'via help options'")
	}
	b, err := via.NewBuilderByName(config, ctx.Args().First())
	if err != nil {
		return err
	}
	c := b.Context.StageDir.Join("configure").String()
	cmd := exec.Command("sh", c, "--help")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func create(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("pack requires a 'URL' argument. see: 'via help create'")
	}
	err := via.Create(config, ctx.Args().First(), "extra")
	if err != nil {
		return err
	}
	return nil
}

func pack(ctx *cli.Context) error {
	via.Verbose(ctx.Bool("v"))
	for _, arg := range ctx.Args().Slice() {
		plan, err := via.NewPlan(config, arg)
		if err != nil {
			return err
		}
		b := via.NewBuilder(config, plan)
		if err := b.Package(b.Context.BuildDir); err != nil {
			return err
		}
	}
	return nil
}

func debug(_ *cli.Context) error {
	cmds := []string{"gcc", "g++", "python", "ld", "perl", "make", "bash", "ccache", "strip"}
	env := config.SanitizeEnv()
	sort.Strings(env)
	for _, v := range env {
		e := strings.SplitN(v, "=", 2)
		console.Println(e[0], e[1])
	}
	console.Flush()
	fmt.Println("PATHS:")
	for _, p := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		console.Println(p)
	}
	for _, c := range cmds {
		fmt.Printf("%s:\n", strings.ToUpper(c))
		execs("which", "-a", c)
		execs(c, "--version")
	}
	return nil
}

// Executes 'cmd' with 'args' useing os.Stdout and os.Stderr
func execs(cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Stderr = os.Stderr
	e.Stdout = os.Stdout
	return e.Run()
}

// Finds all locations of each 'cmd' in PATH and prints to stdout
func which(cmds ...string) {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, c := range cmds {
		fmt.Printf("%s:\n", strings.ToUpper(c))
		for _, p := range paths {
			j := via.Path(p).Join(c)
			if j.Exists() {
				fmt.Println(j)
			}
		}
	}
}

func owns(ctx *cli.Context) error {
	files, err := via.ReadRepoFiles(config)
	if err != nil {
		return err
	}
	for _, arg := range ctx.Args().Slice() {
		owners := files.Owners(arg)
		if len(owners) == 0 {
			console.Println("file:", arg, "owner not found")
			continue
		}
		console.Println("file:", arg, "owners:", owners)
	}
	console.Flush()
	return nil
}

func cd(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("cd requires a 'PLAN' argument. see: 'via help cd'")
	}
	b, err := via.NewBuilderByName(config, ctx.Args().First())
	if err != nil {
		return err
	}
	if ctx.Bool("s") {
		fmt.Printf("cd %s", b.Context.StageDir)
		return nil
	}
	if ctx.Bool("b") {
		fmt.Printf("cd %s", b.Context.BuildDir)
		return nil
	}
	return fmt.Errorf("cd requires either -s or -b flag")
}
