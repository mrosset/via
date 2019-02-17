package main

import (
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"github.com/mrosset/util/file"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	bindRO = syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REC
	bindRW = syscall.MS_BIND | syscall.MS_REC
)

var (
	containCommands = []*cli.Command{
		&cli.Command{
			Name:   "enter",
			Usage:  "enter build namespace",
			Action: contain,
		},
	}
)

func init() {
	var err error
	if config, err = readconfig(); err != nil {
		elog.Fatal(err)
	}
}

func init() {
	reexec.Register("init", initialize)
	if reexec.Init() {
		os.Exit(0)
	}
	app.Commands = append(app.Commands, containCommands...)
}

func linksh(root string) error {
	var (
		source = filepath.Join(config.Prefix, "bin", "bash")
		bin    = filepath.Join(root, "bin")
		target = filepath.Join(bin, "sh")
	)
	if err := os.MkdirAll(bin, 0755); err != nil {
		return err
	}
	return os.Link(source, target)
}

func bindbin(root string) error {
	var (
		source = filepath.Join(config.Prefix, "bin")
		target = filepath.Join(root, "bin")
	)
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	return syscall.Mount(
		source,
		target,
		"",
		bindRO,
		"",
	)
}

// instead of linking, bind sh to bash to avoid cross linking across
// devices
func bindsh(root string) error {
	var (
		source = filepath.Join(config.Prefix, "bin", "bash")
		bin    = filepath.Join(root, "bin")
		target = filepath.Join(bin, "sh")
	)
	if err := os.MkdirAll(bin, 0755); err != nil {
		return err
	}
	if err := file.Touch(target); err != nil {
		return err
	}
	return syscall.Mount(
		source,
		target,
		"",
		bindRO,
		"",
	)
}

func initialize() {
	root, err := ioutil.TempDir("", "via-build")
	if err != nil {
		elog.Fatal(err)
	}
	// set hostname
	if err := syscall.Sethostname([]byte("via-build")); err != nil {
		elog.Fatalf("could not set hostname: %s", err)
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		elog.Fatal(err)
	}
	// setup all our mounts
	if err := mount(root); err != nil {
		elog.Fatal(err)
	}
	// setup busybox and links
	if err := busybox(root); err != nil {
		elog.Fatal(err)
	}
	// finally pivot our root
	if err := pivot(root); err != nil {
		elog.Fatal(err)
	}
	run()
}
func run() {
	var (
		path = "/bin/sh"
		args = []string{}
	)

	switch {
	case len(os.Args) >= 2 && os.Args[1] == "build":
		path = viabin
		args = append([]string{"via"}, os.Args[1:]...)
	case len(os.Args) >= 2 && os.Args[1] == "contain":
		path = "/bin/sh"
		args = []string{}
	default:
		elog.Fatalf("can not handle arguments %v", os.Args)
	}

	cmd := &exec.Cmd{
		Path:   path,
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env: []string{
			fmt.Sprintf("TERM=%s", os.Getenv("TERM")),
			fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
			fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")),
			fmt.Sprintf("CFLAGS=%s", config.Env["CFLAGS"]),
			fmt.Sprintf("LDFLAGS=%s", config.Env["LDFLAGS"]),
			fmt.Sprintf("PATH=%s/bin:/bin:/home/mrosset/gocode/bin", config.Prefix),
			"PS1=-[via-build]- # ",
		},
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
	if err := cmd.Run(); err != nil {
		elog.Fatal(err)
	}
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

type fileSystem struct {
	Source string
	Type   string
	Target string
	Flags  int
	Data   string
	MakeFn func(string) error
}

func (fs fileSystem) Mount(root string) error {
	target := filepath.Join(root, fs.Target)
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	return syscall.Mount(
		fs.Source,
		target,
		fs.Type,
		uintptr(fs.Flags),
		fs.Data,
	)
}

func mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func busybox(root string) error {
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0755); err != nil {
		return err
	}
	bpath := filepath.Join(config.Prefix, "bin", "busybox")
	cmd := exec.Cmd{
		Path:   bpath,
		Args:   []string{"busybox", "--install", "-s", bin},
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	out, err := os.OpenFile(filepath.Join(bin, "busybox"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	return file.Copy(out, bpath)
}

func bind(source, root string) error {
	if source == "" {
		return fmt.Errorf("source can not be ''")
	}
	target := filepath.Join(root, source)
	stat, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("bind %s to %s with error %s", source, target, err)
	}
	if stat.IsDir() {
		os.MkdirAll(target, 0755)
	} else {
		dir := filepath.Dir(target)
		os.MkdirAll(dir, 0755)
		if err := file.Touch(target); err != nil {
			return err
		}
	}
	return syscall.Mount(
		source,
		target,
		"",
		bindRO,
		"",
	)
}

func mount(root string) error {
	// our binds
	binds := []string{
		"/dev",
		"/etc/resolv.conf",
		"/etc/ssl",
		"/etc/passwd",
		os.ExpandEnv("$HOME/.ccache"),
		config.Cache.String(),
		config.Plans.String(),
		config.Repo.String(),
		config.Prefix,
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
	// mount our binds
	for _, source := range binds {
		if err := bind(source, root); err != nil {
			elog.Printf("binding %s to %s", source, filepath.Join(root, source))
			return err
		}
	}
	// mount our filesystems
	for _, m := range fs {
		if err := m.Mount(root); err != nil {
			elog.Printf("mounting %s to %s", m.Source, filepath.Join(root, m.Source))
			return err
		}
	}
	return nil
}

func pivot(newroot string) error {
	oldroot := filepath.Join(newroot, "/.root")

	// bind mount newroot to itself - this is a slight hack
	// needed to work around a pivot_root requirement
	if err := syscall.Mount(
		newroot,
		newroot,
		"",
		bindRO,
		"",
	); err != nil {
		return err
	}

	// create oldroot directory
	if err := os.MkdirAll(oldroot, 0700); err != nil {
		return err
	}

	// call pivot_root
	if err := syscall.PivotRoot(newroot, oldroot); err != nil {
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
	return os.RemoveAll(oldroot)
}
