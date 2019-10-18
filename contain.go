package main

import (
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"github.com/mrosset/util/file"
	"github.com/mrosset/via/pkg"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

const (
	bindRO = syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REC
	bindRW = syscall.MS_BIND | syscall.MS_REC
)

var (
	containCommands = []*cli.Command{
		&cli.Command{
			Name:    "enter",
			Aliases: []string{"shell"},
			Usage:   "enter build namespace",
			Action:  contain,
		},
	}

	defaultPaths = fmt.Sprintf("%s:%s:/opt/via/bin:/bin",
		config.Env["PATH"],
		via.Path(os.Getenv("GOPATH")).Join("bin"),
	)

	defaultEnv = via.Env{
		"INSIDE_VIA": "true",
		"TERM":       os.Getenv("TERM"),
		"HOME":       os.Getenv("HOME"),
		"GOPATH":     os.Getenv("GOPATH"),
		"CLFAGS":     config.Env["CLFAGS"],
		"LDFLAGS":    config.Env["LDFLAGS"],
		"CDPATH":     viapath.Dir().String(),
		"PATH":       os.ExpandEnv(defaultPaths),
		"PS1":        "-[via-build]- $ ",
	}
)

func init() {
	reexec.Register("init", initialize)
	if reexec.Init() {
		os.Exit(0)
	}
	app.Commands = append(app.Commands, containCommands...)
}

func linksh(root string) error {
	var (
		source = config.Prefix.Join("bin", "bash")
		bin    = via.Path(root).Join("bin")
		target = bin.Join("sh")
	)
	if err := bin.MkdirAll(); err != nil {
		return err
	}
	return os.Link(source.String(), target.String())
}

func bindbin(root via.Path) error {
	var (
		source = config.Prefix.Join("bin")
		target = root.Join("bin")
	)
	if err := target.MkdirAll(); err != nil {
		return err
	}
	return syscall.Mount(
		source.String(),
		target.String(),
		"",
		bindRO,
		"",
	)
}

// instead of linking, bind sh to bash to avoid cross linking across
// devices
func bindsh(root string) error {
	var (
		source = config.Prefix.Join("bin", "bash")
		bin    = via.NewPath(root, "bin")
		target = bin.Join("sh")
	)
	if err := bin.MkdirAll(); err != nil {
		return err
	}
	if err := target.Touch(); err != nil {
		return err
	}
	return syscall.Mount(
		source.String(),
		target.String(),
		"",
		bindRO,
		"",
	)
}

func initialize() {
	root := config.Cache.Root()
	config.Root = config.Cache.Root()

	// populate root
	if err := populate(); err != nil {
		elog.Fatal(err)
	}
	// set hostname
	if err := syscall.Sethostname([]byte("via-build")); err != nil {
		elog.Fatalf("could not set hostname: %s", err)
	}
	if err := root.MkdirAll(); err != nil {
		elog.Fatal(err)
	}
	// setup all our mounts
	if err := mount(root); err != nil {
		elog.Fatal(err)
	}
	if err := bindbin(root); err != nil {
		elog.Fatal(err)
	}
	// setup busybox and links
	// if err := busybox(root); err != nil {
	//	elog.Fatal(err)
	// }
	// finally pivot our root
	if err := pivot(root); err != nil {
		elog.Fatal(err)
	}
	if err := enter(root); err != nil {
		elog.Fatal(err)
	}
}

func populate() error {
	if via.IsInstalled(config, "devel") {
		return nil
	}
	devel, err := via.NewPlan(config, "devel")
	if err != nil {
		return err
	}
	fmt.Println("populating namespace")
	b := via.NewBatch(config, ioutil.Discard)
	b.Walk(devel)
	if errors := b.Install(); len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// Enter names space and either runs build or a shell
func enter(root via.Path) error {
	var (
		path     = config.Prefix.Join("bin/bash")
		args     = []string{"--login"}
		profile  = root.Join("etc", "profile")
		complete = []byte("eval $(via --init-completion bash)\n")
	)

	if !profile.Exists() {
		if err := ioutil.WriteFile(profile.String(), complete, 0644); err != nil {
			return err
		}
	}

	switch {
	case len(os.Args) >= 2 && os.Args[1] == "build":
		path = viabin
		args = append([]string{"via"}, os.Args[1:]...)
	case len(os.Args) >= 2 && os.Args[1] == "contain":
	default:
		return fmt.Errorf("can not handle arguments %v", os.Args)
	}

	cmd := &exec.Cmd{
		Path:   path.String(),
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    defaultEnv.KeyValue(),
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUSER,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 1000,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 1001,
					HostID:      os.Getgid(),
					Size:        1,
				},
			},
		},
		Dir: os.Getenv("HOME"),
	}
	return cmd.Run()
}

func contain(ctx *cli.Context) error {
	var (
		args = []string{}
	)
	if len(os.Args) > 1 && os.Args[1] == "build" {
		args = []string{"init", "build", "-real"}
	} else {
		args = []string{"init", "contain"}
	}
	// maybe there is a better way to chain down flags?
	for _, f := range ctx.FlagNames() {
		args = append(args, fmt.Sprintf("-%s", f))
	}
	args = append(args, ctx.Args().Slice()...)
	cmd := reexec.Command(args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	return cmd.Run()
}

func clean() error {
	config.Root = config.Cache.Root()
	names, err := config.DB.InstalledPlans(config)
	if err != nil {
		return err
	}
	for _, n := range names {
		if err := via.Remove(config, n); err != nil {
			return err
		}
	}
	return err
}

type fileSystem struct {
	Source string
	Type   string
	Target string
	Flags  int
	Data   string
	MakeFn func(string) error
}

func (fs fileSystem) Mount(root string) error {
	target := via.Path(root).Join(fs.Target)
	if err := target.MkdirAll(); err != nil {
		return err
	}
	return syscall.Mount(
		fs.Source,
		target.String(),
		fs.Type,
		uintptr(fs.Flags),
		fs.Data,
	)
}

func busybox(root via.Path) error {
	var (
		bin    = via.Path.Join(root, "bin")
		source = ""
		target = bin.Join("busybox")
		err    error
	)
	if source, err = exec.LookPath("busybox"); err != nil {
		return err
	}
	if err := bin.MkdirAll(); err != nil {
		return err
	}

	out, err := os.OpenFile(target.String(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	if err := file.Copy(out, source); err != nil {
		return err
	}
	out.Close()
	cmd := exec.Cmd{
		Path:   target.String(),
		Args:   []string{"busybox", "--install", bin.String()},
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
	return cmd.Run()
}

func bind(source, root via.Path) error {
	if source == "" {
		return fmt.Errorf("source can not be ''")
	}
	var (
		target = root.Join(source.String())
	)
	stat, err := source.Stat()
	if err != nil {
		return fmt.Errorf("bind %s to %s with error %s", source, target, err)
	}
	if stat.IsDir() {
		target.Ensure()
	} else {
		if err := target.Dir().MkdirAll(); err != nil {
			return err
		}
		if err := target.Touch(); err != nil {
			return err
		}
	}
	return syscall.Mount(
		source.String(),
		target.String(),
		"",
		bindRO,
		"",
	)
}

func mount(root via.Path) error {
	// our binds
	binds := []via.Path{
		"/dev", // TODO: don not bind all of dev. mknod's?
		"/etc/resolv.conf",
		"/etc/ssl",
		"/etc/passwd",
		"/etc/group",
		via.Path("$HOME/.ccache").Expand(),
		config.Cache.ToPath(),
		config.Plans.ToPath(),
		config.Repo.ToPath(),
		viabin,
	}
	// our filesystems
	fs := []fileSystem{
		{
			Source: "proc",
			Target: "/proc",
			Type:   "proc",
		},
		{
			Source: "tmpfs",
			Target: "/tmp",
			Type:   "tmpfs",
		},
	}
	// mount our filesystems
	for _, m := range fs {
		if err := m.Mount(root.String()); err != nil {
			return err
		}
	}
	// mount our binds
	for _, source := range binds {
		if err := bind(source, root); err != nil {
			return err
		}
	}
	return nil
}

func pivot(newroot via.Path) error {
	oldroot := newroot.Join("/.root")

	// bind mount newroot to itself - this is a slight hack
	// needed to work around a pivot_root requirement
	if err := syscall.Mount(
		newroot.String(),
		newroot.String(),
		"",
		bindRO,
		"",
	); err != nil {
		return err
	}

	// create oldroot directory
	if err := os.MkdirAll(oldroot.String(), 0700); err != nil {
		return err
	}

	// call pivot_root
	if err := syscall.PivotRoot(newroot.String(), oldroot.String()); err != nil {
		return err
	}

	// ensure current working directory is set to new root
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount oldroot, which now lives at /.pivot_root
	if err := syscall.Unmount("/.root", syscall.MNT_DETACH); err != nil {
		return err
	}
	return oldroot.RemoveAll()
}
