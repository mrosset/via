package main

import (
	"fmt"
	"github.com/mrosset/gurl"
	"github.com/mrosset/util/console"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

var (
	elog   = log.New(os.Stderr, "", log.Lshortfile)
	lfmt   = "%-20.20s %v\n"
	config = via.GetConfig()
	app    = &cli.App{
		Name:  "via",
		Usage: "a systems package manager",
	}

	// build command
	cbuild = &cli.Command{
		Name:   "build",
		Usage:  "builds a plan locally",
		Action: local,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "c",
				Value: false,
				Usage: "clean build directory before building",
			},
			&cli.BoolFlag{
				Name:  "v",
				Value: false,
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
				Value: true,
				Usage: "install package after building",
			},
			&cli.BoolFlag{
				Name:  "u",
				Value: false,
				Usage: "force downloading of sources",
			},
			&cli.BoolFlag{
				Name:  "r",
				Value: false,
				Usage: "builds plan using daemon",
			},
		},
	}

	// install command
	cinstall = &cli.Command{
		Name:   "install",
		Usage:  "installs package",
		Action: batch,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "r",
				Value: config.Root,
				Usage: "use `\"DIR\"` as root",
			},
			&cli.BoolFlag{
				Name:  "y",
				Value: false,
				Usage: "Don't prompt to install",
			},
		},
	}

	// remove command
	cremove = &cli.Command{
		Name:   "remove",
		Usage:  "uninstall package",
		Action: remove,
	}

	// edit command
	cedit = &cli.Command{
		Name:   "edit",
		Usage:  "calls EDITOR to edit plan",
		Action: edit,
	}

	// show command
	cshow = &cli.Command{
		Name:   "show",
		Usage:  "prints plan to stdout",
		Action: show,
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
		},
	}

	// config command
	cconfig = &cli.Command{
		Name:   "config",
		Usage:  "prints config to stdout",
		Action: fconfig,
	}

	// list command
	clist = &cli.Command{
		Name:   "list",
		Usage:  "list files for `PLAN`",
		Action: list,
	}

	// lint command
	clint = &cli.Command{
		Name:   "lint",
		Usage:  "lint and format plans",
		Action: lint,
	}

	// repo command
	crepo = &cli.Command{
		Name:   "repo",
		Usage:  "recreates file db",
		Action: repo,
	}

	clog = &cli.Command{
		Name:   "log",
		Usage:  "output's config.log for build",
		Action: plog,
	}

	celf = &cli.Command{
		Name:   "elf",
		Usage:  "prints elf information to stdout",
		Action: elf,
	}

	cdiff = &cli.Command{
		Name:   "diff",
		Usage:  "diff's plan working directory against git HEAD",
		Action: diff,
	}

	csearch = &cli.Command{
		Name:   "search",
		Usage:  "lists all of the available packages",
		Action: search,
	}

	coptions = &cli.Command{
		Name:   "options",
		Usage:  "prints the GNU configure options for a package",
		Action: options,
	}

	cstrap = &cli.Command{
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
	}

	ccreate = &cli.Command{
		Name:   "create",
		Usage:  "creates a plan from tarball URL",
		Action: create,
	}

	cdaemon = &cli.Command{
		Name:   "daemon",
		Usage:  "starts build daemon",
		Action: daemon,
	}

	chash = &cli.Command{
		Name:   "hash",
		Usage:  "DEV ONLY sync the plans Oid with binary banch",
		Action: hash,
	}

	cpack = &cli.Command{
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
	}

	cdebug = &cli.Command{
		Name:   "debug",
		Usage:  "displays enviroment and PATH details",
		Action: debug,
	}

	cowns = &cli.Command{
		Name:   "owns",
		Usage:  "find which plans owns 'file'",
		Action: owns,
	}

	cfix = &cli.Command{
		Name:   "fix",
		Usage:  "DEV ONLY used to mass modify plans",
		Action: fix,
	}

	cclean = &cli.Command{
		Name:   "clean",
		Usage:  "cleans cache directory",
		Action: clean,
	}

	ccd = &cli.Command{
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
	}

	cget = &cli.Command{
		Name:   "get",
		Usage:  "gets 'plans' sources from ipfs into current directory",
		Action: get,
	}

	// cadd = &cli.Command{
	//	Name:   "add",
	//	Usage:  "Adds 'dir' to ipfs and saves plan SourceCid",
	//	Action: add,
	//	Flags: []cli.Flag{
	//		&cli.StringFlag{
	//			Name:  "p",
	//			Value: "",
	//			Usage: "plan to add source directory to",
	//		},
	//	},
	// }
)

func main() {
	app.Commands = []*cli.Command{
		cinstall,
		cremove,
		cbuild,
		clist,
		cconfig,
		cshow,
		crepo,
		clint,
		cedit,
		clog,
		celf,
		cdiff,
		csearch,
		coptions,
		cstrap,
		ccreate,
		cdaemon,
		chash,
		cpack,
		cdebug,
		cowns,
		cfix,
		cclean,
		ccd,
		cget,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// func add(ctx *cli.Context) error {
//	if !ctx.Args().Present() {
//		return fmt.Errorf("get requires a 'dir' argument. see: 'via help add'")
//	}
//	if ctx.String("p") == "" {
//		return fmt.Errorf("add requires -p 'plan' flag. see: 'via help add'")
//	}
//	plan, err := via.NewPlan(ctx.String("p"))
//	if err != nil {
//		return err
//	}
//	cid, err := via.IpfsAdd(via.Path(ctx.Args().First()), false)
//	if err != nil {
//		return err
//	}
//	plan.SourceCid = cid
//	return plan.Save()
// }

func get(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("get requires a 'PLAN' argument. see: 'via help get'")
	}

	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}
	return gurl.Download("./", plan.Expand().Url)
}

func clean(ctx *cli.Context) error {
	if err := os.RemoveAll(via.Path(config.Cache.Builds()).String()); err != nil {
		return err
	}
	if err := os.RemoveAll(via.Path(config.Cache.Stages()).String()); err != nil {
		return err
	}
	return nil
}

func fix(ctx *cli.Context) error {
	plans, err := via.GetPlans()
	if err != nil {
		return err
	}
	for _, p := range plans {
		p.IsRebuilt = false
		p.Save()
	}
	return nil
}

func daemon(ctx *cli.Context) error {
	return via.StartDaemon()
}

func strap(ctx *cli.Context) error {

	dplan, err := via.NewPlan("emacs")

	if err != nil {
		return err
	}

	for _, p := range dplan.ManualDepends {
		plan, err := via.NewPlan(p)
		if err != nil {
			return err
		}
		if ctx.Bool("m") {
			plan.IsRebuilt = false
			plan.Save()
			continue
		}
		if plan.IsRebuilt {
			fmt.Printf(lfmt, "rebuilt", plan.NameVersion())
			continue
		}
		via.Clean(plan.Name)

		if err := via.BuildSteps(plan); err != nil {
			return err
		}

		plan.IsRebuilt = true
		plan.Save()
	}
	return nil
}

func batch(ctx *cli.Context) error {
	var errors []error
	if !ctx.Args().Present() {
		return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
	}

	via.Root(ctx.String("r"))
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}
	batch := via.NewBatch(config)
	batch.Walk(plan)

	switch ctx.Bool("y") {
	case false:
		errors = batch.PromptInstall()
	case true:
		errors = batch.Install()

	}
	if len(errors) > 0 {
		log.Fatal(errors)
	}
	return nil
}

func install(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
	}

	via.Root(ctx.String("r"))

	for _, arg := range ctx.Args().Slice() {
		if err := via.Install(arg); err != nil {
			return err
		}
	}
	return nil
}

func remove(ctx *cli.Context) error {
	return via.Remove(ctx.Args().First())
}

func local(ctx *cli.Context) error {
	// if r flag build package with RPC daemon
	if ctx.Bool("r") {
		return remote(ctx)
	}
	if !ctx.Args().Present() {
		return fmt.Errorf("build requires a 'PLAN' argument. see: 'via help build'")
	}
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}

	via.Verbose(ctx.Bool("v"))
	via.Debug(ctx.Bool("d"))
	via.Update(ctx.Bool("u"))

	if ctx.Bool("c") {
		via.Clean(plan.Name)
	}
	if ctx.Bool("dd") {
		err = via.BuildDeps(plan)
		if err != nil {
			return err
		}
	} else {
		err = via.BuildSteps(plan)
		if err != nil {
			return err
		}
	}
	if ctx.Bool("i") {
		return via.Install(plan.Name)
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
	p, _ := via.NewPlan(ctx.Args().First())
	req := via.Request{*p}
	return c.Call("Builder.RpcBuild", req, &res)
}

func edit(ctx *cli.Context) error {
	var (
		editor = os.Getenv("EDITOR")
		arg0   = ctx.Args().First()
		p      = filepath.Join(config.Plans, "config.json")
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
	err = cmd.Run()
	if err != nil {
		return err
	}
	elog.Println("linting...")
	via.Verbose(false)
	return via.Lint()
}

func list(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("list requires a 'PLAN' argument. see: 'via help list'")
	}
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}
	for _, f := range plan.Files {
		fmt.Println(f)
	}
	return nil
}

func lint(ctx *cli.Context) error {
	return via.Lint()
}

func show(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("show requires a 'PLAN' argument. see: 'via help show'")
	}
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		elog.Fatal(err)
	}
	if ctx.String("t") != "" {
		tmpl, err := template.New("stdout").Parse(ctx.String("t") + "\n")
		if err != nil {
			panic(err)
		}
		return tmpl.Execute(os.Stdout, plan)
	}
	if ctx.Bool("d") {
		tmpl, err := template.New("stdout").Parse("{{.AutoDepends}}\n")
		if err != nil {
			panic(err)
		}
		return tmpl.Execute(os.Stdout, plan)
	}
	err = json.WritePretty(&plan, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func fconfig(ctx *cli.Context) error {
	err := json.WritePretty(via.GetConfig(), os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func repo(ctx *cli.Context) error {
	return via.RepoCreate()
}

func plog(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("show requires a 'PLAN' argument. see: 'via help log'")
	}
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}
	f := filepath.Join(plan.BuildDir(), "config.log")
	err = file.Cat(os.Stdout, f)
	if err != nil {
		log.Fatal(err)
	}
	return nil
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
		glob := filepath.Join(config.Plans, "*", arg+".json")
		res, err := filepath.Glob(glob)
		if err != nil {
			return err
		}
		git := exec.Command("git", "diff", strings.Join(res, " "))
		git.Dir = config.Plans
		git.Stdout = os.Stdout
		git.Stderr = os.Stderr
		err = git.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func search(ctx *cli.Context) error {
	plans, err := via.GetPlans()
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
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}
	c := filepath.Join(plan.GetStageDir(), "configure")
	fmt.Println(c)
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
	err := via.Create(ctx.Args().First(), "core")
	if err != nil {
		return err
	}
	return nil
}

func hash(ctx *cli.Context) error {
	via.SyncHashs()
	return nil
}

func pack(ctx *cli.Context) error {
	via.Verbose(ctx.Bool("v"))
	for _, arg := range ctx.Args().Slice() {
		plan, err := via.NewPlan(arg)
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

func debug(ctx *cli.Context) error {
	cmds := []string{"gcc", "g++", "python", "make", "bash", "ld"}
	env := os.Environ()
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

func wversion(path string) {
}

// Finds all locations of each 'cmd' in PATH and prints to stdout
func which(cmds ...string) {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, c := range cmds {
		fmt.Printf("%s:\n", strings.ToUpper(c))
		for _, p := range paths {
			j := filepath.Join(p, c)
			if file.Exists(j) {
				fmt.Println(j)
			}
		}
	}
}

func owns(ctx *cli.Context) error {
	rfiles, err := via.ReadRepoFiles()
	if err != nil {
		return err
	}
	for _, arg := range ctx.Args().Slice() {
		owner := rfiles.Owns(arg)
		if owner == "" {
			fmt.Println(arg+":", "owner not found.")
			continue
		}
		fmt.Println(owner)
	}
	return nil
}

func cd(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("cd requires a 'PLAN' argument. see: 'via help cd'")
	}
	plan, err := via.NewPlan(ctx.Args().First())
	if err != nil {
		return err
	}
	if ctx.Bool("s") {
		fmt.Printf("cd %s", plan.GetStageDir())
		return nil
	}
	if ctx.Bool("b") {
		fmt.Printf("cd %s", plan.BuildDir())
		return nil
	}
	return fmt.Errorf("cd requires either -s or -b flag")
}

/*
func add() error {
	if len(command.Args()) < 1 {
		return errors.New("no plans specified")
	}
	for _, arg := range command.Args() {
		glob := filepath.Join(config.Plans, "*", arg+".json")
		res, err := filepath.Glob(glob)
		if err != nil {
			return err
		}
		git := exec.Command("git", "add", strings.Join(res, " "))
		git.Dir = config.Plans
		git.Stdout = os.Stdout
		git.Stderr = os.Stderr
		err = git.Run()
		if err != nil {
			return err
		}
	}
*/

/*

	return nil
}

func checkout() error {
	if len(command.Args()) < 1 {
		return errors.New("git branch needs to be specified")
	}
	arg := command.Args()[0]
	git := exec.Command("git", "checkout", arg)
	git.Dir = config.Plans
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	return git.Run()
}

func branch() error {
	git := exec.Command("git", "branch")
	git.Dir = config.Plans
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	return git.Run()

}

func sync() error {
	return via.PlanSync()
}

func oldCommands() {
	// Old Flags
	root     = flag.String("r", "/", "root directory")
	verbose  = flag.Bool("v", false, "verbose output")
	finstall = flag.Bool("i", true, "install package after build (default true)")
	fdebug   = flag.Bool("d", false, "debug output")
	fclean   = flag.Bool("c", false, "clean before build")
	fupdate  = flag.Bool("u", false, "force download source")
	fdeps    = flag.Bool("deps", false, "build depends if needed")

	// Old Commands
	flag.Parse()
	via.Verbose(*verbose)
	via.Update(*fupdate)
	via.Deps(*fdeps)

	via.Root(*root)
	util.Verbose = *verbose
	via.Debug(*fdebug)
	command.Add("add", add, "add plan/s to git index")
	command.Add("branch", branch, "prints plan branch to stdout")
	command.Add("cd", cd, "returns a bash evaluable cd path")
	command.Add("checkout", checkout, "changes plan branch")
	command.Add("clean", clean, "clean build dir")
	command.Add("create", create, "create plan from URL")
	command.Add("diff", diff, "prints git diff for plan(s)")
	command.Add("elf", elf, "prints elf information to stdout")
	command.Add("ipfs", ipfs, "test ipfs connection")
	command.Add("lint", lint, "lint plans")
	command.Add("log", plog, "print config log for plan")
	command.Add("owns", owns, "finds which package owns a file")
	command.Add("options", options, "prints the GNU configure options for a package")
	command.Add("pack", pack, "package plan")
	command.Add("remove", remove, "remove package")
	command.Add("repo", repo, "update repo")
	command.Add("search", search, "search for plans (currently lists all use grep)")
	command.Add("sync", sync, "fetch remote repo data")
	command.Add("synchashs", synchashs, "DEV ONLY sync the plans Oid with binary banch")
	if *fdebug {
		pdebug()
	}
	err = command.Run()
	if err != nil {
		elog.Fatal(err)
	}
	return
}
*/
