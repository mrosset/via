package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/mrosset/util"
	"github.com/mrosset/util/console/command"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	elog     = log.New(os.Stderr, "", log.Lshortfile)
	root     = flag.String("r", "/", "root directory")
	verbose  = flag.Bool("v", false, "verbose output")
	finstall = flag.Bool("i", true, "install package after build (default true)")
	fdebug   = flag.Bool("d", false, "debug output")
	config   = via.GetConfig()
	fclean   = flag.Bool("c", false, "clean before build")
	fupdate  = flag.Bool("u", false, "force download source")
	fdeps    = flag.Bool("deps", false, "build depends if needed")
)

func main() {
	flag.Parse()
	via.Verbose(*verbose)
	via.Update(*fupdate)
	via.Deps(*fdeps)

	via.Root(*root)
	util.Verbose = *verbose
	via.Debug(*fdebug)
	command.Add("add", add, "add plan/s to git index")
	command.Add("branch", branch, "prints plan branch to stdout")
	command.Add("build", build, "build plan")
	command.Add("dock", dock, "build plan inside docker")
	command.Add("start", start, "starts rpc server")
	command.Add("cd", cd, "returns a bash evaluable cd path")
	command.Add("checkout", checkout, "changes plan branch")
	command.Add("clean", clean, "clean build dir")
	command.Add("config", fnConfig, "prints config to stdout")
	command.Add("create", create, "create plan from URL")
	command.Add("diff", diff, "prints git diff for plan(s)")
	command.Add("edit", edit, "calls EDITOR to edit plan")
	command.Add("elf", elf, "prints elf information to stdout")
	command.Add("install", install, "install package")
	command.Add("ipfs", ipfs, "test ipfs connection")
	command.Add("lint", lint, "lint plans")
	command.Add("list", list, "lists files")
	command.Add("log", plog, "print config log for plan")
	command.Add("owns", owns, "finds which package owns a file")
	command.Add("options", options, "prints the GNU configure options for a package")
	command.Add("pack", pack, "package plan")
	command.Add("remove", remove, "remove package")
	command.Add("repo", repo, "update repo")
	command.Add("search", search, "search for plans (currently lists all use grep)")
	command.Add("show", fnShow, "prints plan to stdout")
	command.Add("sync", sync, "fetch remote repo data")
	command.Add("synchashs", synchashs, "DEV ONLY sync the plans Oid with binary banch")
	if *fdebug {
		pdebug()
	}
	err := command.Run()
	if err != nil {
		elog.Fatal(err)
	}
	return
}

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
func edit() error {
	var (
		editor = os.Getenv("EDITOR")
		arg0   = command.Args()[0]
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

func create() error {
	url := command.Args()[0]
	group := command.Args()[1]
	err := via.Create(url, group)
	if err != nil {
		return err
	}
	return nil
}

func plog() error {
	for _, arg := range command.Args() {
		plan, err := via.NewPlan(arg)
		if err != nil {
			return err
		}
		f := filepath.Join(plan.BuildDir(), "config.log")
		err = file.Cat(os.Stdout, f)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func start() error {
	if err := via.Install("devel"); err != nil {
		return err
	}
	pdebug()
	return via.Listen()
}
func dock() error {
	arg := command.Args()[0]
	c, dc := via.NewRpcClient()
	c.Start()
	defer c.Stop()
	_, err := dc.Call("build", arg)
	return err
}

func odock() error {
	// FIXME this breaks when no argument is passed
	arg := command.Args()[0]
	dargs := []string{
		"run", "-it",
		"-v", "/home:/home",
		"strings/via:devel",
		"/usr/bin/via",
		"-c",
		"-deps",
		"build",
		arg,
	}
	cmd := exec.Command("docker", dargs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func build() error {
	for _, arg := range command.Args() {
		arg = strings.Replace(arg, ".json", "", 1)
		if *fclean {
			via.Clean(arg)
		}
		plan, err := via.NewPlan(arg)
		if err != nil {
			return err
		}
		if *fdeps {
			err := via.BuildDeps(plan)
			if err != nil {
				return err
			}
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

func pack() error {
	for _, arg := range command.Args() {
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

func list() error {
	for _, arg := range command.Args() {
		plan, err := via.NewPlan(arg)
		if err != nil {
			return err
		}
		for _, f := range plan.Files {
			fmt.Println(f)
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

func fnShow() error {
	for _, arg := range command.Args() {
		plan, err := via.NewPlan(arg)
		if err != nil {
			elog.Fatal(err)
		}
		buf := new(bytes.Buffer)
		less := exec.Command("less")
		less.Stdin = buf
		less.Stdout = os.Stdout
		less.Stderr = os.Stderr
		err = json.WritePretty(&plan, buf)
		if err != nil {
			fmt.Println(err)
		}
		less.Run()
	}
	return nil
}

func fnConfig() error {
	err := json.WritePretty(via.GetConfig(), os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func clean() error {
	return command.ArgsDo(via.Clean)
}

func elf() error {
	for _, arg := range command.Args() {
		err := via.Readelf(arg)
		if err != nil {
			return err
		}
	}
	return nil
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

func options() error {
	if len(command.Args()) < 1 {

	}
	arg := command.Args()[0]

	plan, err := via.NewPlan(arg)
	if err != nil {
		return err
	}
	c := filepath.Join(plan.GetStageDir(), "configure")
	fmt.Println(c)
	cmd := exec.Command(c, "--help")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func repo() error {
	return via.RepoCreate()
}

func search() error {
	plans, err := via.GetPlans()
	if err != nil {
		return err
	}
	plans.SortSize().Print()
	return nil
}
