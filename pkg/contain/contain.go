package contain

import (
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"github.com/mrosset/util/file"
	"github.com/mrosset/via/pkg"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	BIND_RO = syscall.MS_BIND | syscall.MS_RDONLY | syscall.MS_REC
	BIND_RW = syscall.MS_BIND | syscall.MS_REC
)

var (
	elog   = log.New(os.Stderr, "error: ", log.Lshortfile)
	config = via.GetConfig()
)

func init() {
	reexec.Register("init", initialize)
	if reexec.Init() {
		os.Exit(0)
	}

}

func Append(app *cli.App) {
	cmd := &cli.Command{
		Name:   "contain",
		Action: contain,
		Hidden: false,
	}
	app.Commands = append(app.Commands, cmd)
}

func initialize() {
	fmt.Println(os.Args)
	root, err := ioutil.TempDir("", "via-build")
	if err != nil {
		elog.Fatal(err)
	}

	if err := os.MkdirAll(root, 0700); err != nil {
		elog.Fatal(err)
	}
	// setup busybox
	// if err := busybox(root); err != nil {
	//	elog.Fatal(err)
	// }

	// setup all our mounts
	if err := mount(root); err != nil {
		elog.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0755); err != nil {
		elog.Fatal(err)
	}
	if err := os.Link(
		filepath.Join(config.Prefix, "bin", "sh"),
		filepath.Join(root, "bin", "sh"),
	); err != nil {
		elog.Fatal(err)
	}
	// finally pivot our root
	if err := pivot(root); err != nil {
		elog.Fatal(err)
	}
	run()
}

func run() {

	cmd := exec.Command(filepath.Join(config.Prefix, "bin", "bash"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{
		"HOME=/home/mrosset",
		"GOPATH=/home/mrosset/gocode",
		"PATH=/bin:/opt/via/bin:/home/mrosset/gocode/bin",
		"PS1=-[via-build]- # ",
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER,
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
	}

	if err := cmd.Run(); err != nil {
		elog.Fatal(err)
	}
}

func contain(ctx *cli.Context) error {
	cmd := reexec.Command("init", ctx.Command.Name)
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
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

type FileSystem struct {
	Source string
	Type   string
	Flags  int
	Data   string
	MakeFn func(string) error
}

func (fs FileSystem) Mount(root string) error {

	target := filepath.Join(root, fs.Source)
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
	bpath := "/bin/busybox"
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
	target := filepath.Join(root, source)
	stat, err := os.Stat(source)
	if err != nil {
		return err
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
		BIND_RO,
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
		config.Plans,
		config.Repo,
		filepath.Join(os.Getenv("GOPATH"), "bin/via"),
		config.Prefix,
	}
	// our filesystems
	fs := []FileSystem{
		{
			Source: "proc",
			Type:   "proc",
		},
		{
			Source: "tmpfs",
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
		BIND_RO,
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
	oldroot = "/.root"
	if err := syscall.Unmount(
		oldroot,
		syscall.MNT_DETACH,
	); err != nil {
		return err
	}

	return os.RemoveAll(oldroot)
}
