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
	elog = log.New(os.Stderr, "", log.Lshortfile)
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

// install devel group
func devel(root string) error {
	via.Root(root)
	config := via.GetConfig()
	batch := via.NewBatch(config)
	p, err := via.NewPlan(config, "devel")
	if err != nil {
		return err
	}
	batch.Walk(p)
	errors := batch.Install()
	if len(errors) != 0 {
		return errors[0]
	}
	return nil
}
func initialize() {
	fmt.Println(os.Args)
	root, err := ioutil.TempDir("", "via-build")
	if err != nil {
		elog.Fatal(err)
	}

	// if err := devel(root); err != nil {
	//	elog.Fatal(err)
	// }
	// setup busybox
	// if err := busybox(root); err != nil {
	//	elog.Fatal(err)
	// }
	// Setup all our mounts
	if err := mount(root); err != nil {
		elog.Fatal(err)
	}
	if err := os.Link("/opt/via/bin/sh", filepath.Join(root, "bin", "sh")); err != nil {
		elog.Fatal(err)
	}

	// finally pivot our rout
	if err := pivot(root); err != nil {
		elog.Fatal(err)
	}
	run()
}

func run() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{
		"PATH=/bin:/opt/via/bin",
		"GOPATH=/home/mrosset/gocode",
		"HOME=/home/mrosset",
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
				ContainerID: 1000,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
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
	Target string
	Type   string
	Flags  int
	Data   string
	MakeFn func(string) error
}

func (fs FileSystem) Mount() error {
	if fs.MakeFn != nil {
		if err := fs.MakeFn(fs.Target); err != nil {
			return err
		}
	}
	return syscall.Mount(
		fs.Source,
		fs.Target,
		fs.Type,
		uintptr(fs.Flags),
		fs.Data,
	)
}

func mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func bindFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return file.Touch(path)
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

func mount(root string) error {
	// viabin := filepath.Join(os.Getenv("GOPATH"), "bin", "via")
	mounts := []FileSystem{
		{
			Source: "proc",
			Target: filepath.Join(root, "proc"),
			Type:   "proc",
			MakeFn: mkdir,
		},
		{
			Source: "tmpfs",
			Target: filepath.Join(root, "tmp"),
			Type:   "tmpfs",
			MakeFn: mkdir,
		},
		// FIXME: we don't need everything in dev
		{
			Source: "/dev",
			Target: filepath.Join(root, "dev"),
			Flags:  BIND_RO,
			MakeFn: mkdir,
		},
		{
			Source: filepath.Join(os.Getenv("GOPATH"), "bin", "via"),
			Target: filepath.Join(root, "bin", "via"),
			Flags:  BIND_RO,
			MakeFn: bindFile,
		},
		{
			Source: "/etc/ssl",
			Target: filepath.Join(root, "etc", "ssl"),
			Flags:  BIND_RO,
			MakeFn: mkdir,
		},
		{
			Source: "/home/mrosset/.ccache",
			Target: filepath.Join(root, "/home/mrosset/.ccache"),
			Flags:  BIND_RO,
			MakeFn: mkdir,
		},

		{
			Source: "/home/mrosset/.cache/via",
			Target: filepath.Join(root, "/home/mrosset/.cache/via"),
			Flags:  BIND_RO,
			MakeFn: mkdir,
		},
		{
			Source: "/home/mrosset/gocode/src/github.com/mrosset/via",
			Target: filepath.Join(root, "/home/mrosset/gocode/src/github.com/mrosset/via"),
			Flags:  BIND_RO,
			MakeFn: mkdir,
		},
		{
			Source: "/etc/resolv.conf",
			Target: filepath.Join(root, "etc", "resolv.conf"),
			Flags:  BIND_RO,
			MakeFn: bindFile,
		},
		{
			Source: "/opt/via",
			Target: filepath.Join(root, "/opt/via"),
			Flags:  BIND_RO,
			MakeFn: mkdir,
		},
	}
	errors := []error{}
	for _, m := range mounts {
		if err := m.Mount(); err != nil {
			elog.Printf("%+v", m)
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return errors[0]
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
