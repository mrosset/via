package main

import (
	"fmt"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
		},
	}

	// install command
	cinstall = &cli.Command{
		Name:   "install",
		Usage:  "installs package",
		Action: install,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "r",
				Value: config.Root,
				Usage: "use `\"DIR\"` as root",
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
		Usage:  "output's configure.log for build",
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

	cpatch = &cli.Command{
		Name:   "patch",
		Usage:  "patches dynamic linker",
		Action: patch,
	}

	cdaemon = &cli.Command{
		Name:   "daemon",
		Usage:  "starts build daemon",
		Action: daemon,
	}
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
		cpatch,
		cdaemon,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func daemon(ctx *cli.Context) error {
	return via.StartDaemon()
}
func patch(ctx *cli.Context) error {
	fnWalk := func(path string, fi os.FileInfo, err error) error {
		patch := exec.Command(
			"patchelf",
			"--set-rpath", filepath.Join(config.Root, config.Prefix, "lib"),
			"--set-interpreter", filepath.Join(config.Root, config.Prefix, "lib/ld-linux-x86-64.so.2"),
			path,
		)
		patch.Stdout = os.Stdout
		patch.Run()
		return nil
	}
	path := filepath.Join(config.Root, config.Prefix)
	return filepath.Walk(path, fnWalk)
}

func strap(ctx *cli.Context) error {

	dplan, err := via.FindPlan(via.GetConfig(), "devel")

	if err != nil {
		return err
	}

	c, err := via.Connect()
	if err != nil {
		return err
	}

	for _, p := range dplan.ManualDepends {
		plan, err := via.FindPlan(via.GetConfig(), p)
		if ctx.Bool("m") {
			plan.IsRebuilt = false
			plan.Save(config)
			continue
		}
		if plan.IsRebuilt {
			fmt.Printf(lfmt, "rebuilt", plan.NameVersion())
			continue
		}
		if err != nil {
			return err
		}
		via.Clean(plan.Name)

		res := via.Response{}
		req := via.Request{*plan}

		err = c.Call("Builder.RpcBuild", req, &res)

		if err != nil {
			return err
		}
		plan.IsRebuilt = true
		plan.Save(config)
	}
	return nil
}

func install(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("install requires a 'PLAN' argument. see: 'via help install'")
	}

	via.Root(ctx.String("r"))
	return via.Install(ctx.Args().First())
}

func remove(ctx *cli.Context) error {
	return via.Remove(ctx.Args().First())
}

func local(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("build requires a 'PLAN' argument. see: 'via help build'")
	}
	plan, err := via.FindPlan(via.GetConfig(), ctx.Args().First())
	if err != nil {
		return err
	}

	via.Verbose(ctx.Bool("v"))
	via.Debug(ctx.Bool("d"))
	via.Update(ctx.Bool("u"))

	if ctx.Bool("c") {
		via.Clean(plan.Name)
	}
	err = via.BuildSteps(plan)
	if err != nil {
		return err
	}
	if ctx.Bool("i") {
		return via.Install(plan.Name)
	}
	return nil
}

func build(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("build requires a 'PLAN' argument. see: 'via help build'")
	}
	c, err := via.Connect()
	if err != nil {
		return err
	}
	res := via.Response{}
	p, _ := via.FindPlan(via.GetConfig(), ctx.Args().First())
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
		p, err = via.FindPlanPath(via.GetConfig(), arg0)
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
	plan, err := via.FindPlan(via.GetConfig(), ctx.Args().First())
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
	plan, err := via.FindPlan(via.GetConfig(), ctx.Args().First())
	if err != nil {
		elog.Fatal(err)
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
	plan, err := via.FindPlan(via.GetConfig(), ctx.Args().First())
	con := via.NewConstruct(via.GetConfig(), plan)
	if err != nil {
		return err
	}
	f := filepath.Join(con.BuildPath(), "config.log")
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
	git := exec.Command("git", "diff")
	git.Dir = config.Plans
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	return git.Run()
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
		return fmt.Errorf("show requires a 'PLAN' argument. see: 'via help options'")
	}
	plan, err := via.FindPlan(via.GetConfig(), ctx.Args().First())
	con := via.NewConstruct(via.GetConfig(), plan)
	if err != nil {
		return err
	}
	c := filepath.Join(con.PlanStagePath(), "configure")
	fmt.Println(c)
	cmd := exec.Command(c, "--help")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func create(ctx *cli.Context) error {
	if !ctx.Args().Present() {
		return fmt.Errorf("create requires a 'PLAN' argument. see: 'via help options'")
	}
	err := via.Create(ctx.Args().First(), "core")
	if err != nil {
		return err
	}
	return nil
}

/*
func pdebug() {
	path, _ := os.LookupEnv("PATH")
	home, _ := os.LookupEnv("HOME")
	fmt.Println("PATH", path)
	fmt.Println("HOME", home)
	which("GCC", "gcc")
	which("PYTHON", "python")
}

func which(label, path string) {
	fmt.Printf("%s ", label)
	cmd := exec.Command("which", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		elog.Println(err)
	}
}

func ipfs() error {
	res, err := http.Get(config.Binary)
	io.Copy(os.Stdout, res.Body)
	return err
}

func cd() error {
	if len(command.Args()) < 1 {
		return errors.New("you need to specify a config path")
	}
	arg := command.Args()[0]
	switch arg {
	case "plans":
		fmt.Printf("cd %s", config.Plans)
	default:
		err := fmt.Sprintf("config path %s not found", arg)
		return errors.New(err)
	}
	return nil
}
*/

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
	return nil
}
*/

/*
func diff() error {
	if len(command.Args()) < 1 {
		return errors.New("no plans specified")
	}
	for _, arg := range command.Args() {
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

func pack() error {
	for _, arg := range command.Args() {
		plan, err := via.FindPlan(arg)
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

*/

/*
func clean() error {
	return command.ArgsDo(via.Clean)
}

func sync() error {
	return via.PlanSync()
}

func synchashs() error {
	via.SyncHashs()
	return nil
}
func owns() error {
	rfiles, err := via.ReadRepoFiles()
	if err != nil {
		return err
	}
	for _, arg := range command.Args() {
		owner := rfiles.Owns(arg)
		if owner == "" {
			fmt.Println(arg+":", "owner not found.")
			continue
		}
		fmt.Println(owner)
	}
	return nil
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
