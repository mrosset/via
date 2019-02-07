package main

import (
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		log.Println("INIT")
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			log.Fatal(err)
		}
		panic("--this line should have never been executed, congratulations--")
	}
}

func main() {
	factory, err := libcontainer.New("./containers", libcontainer.RootlessCgroupfs, libcontainer.InitArgs(os.Args[0], "init"))

	if err != nil {
		log.Fatal(err)
		return
	}

	abs, err := filepath.Abs("rootfs")
	if err != nil {
		log.Fatal(err)
	}

	config.Rootfs = abs
	container, err := factory.Create("container-id", config)
	if err != nil {
		log.Fatal(err)
		return
	}

	process := &libcontainer.Process{
		Args:            []string{"/bin/sh"},
		Env:             []string{"PATH=/bin"},
		User:            "root",
		Stdin:           os.Stdin,
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		NoNewPrivileges: func() *bool { b := true; return &b }(),
	}

	// Fails to start process here
	if err := container.Run(process); err != nil {
		container.Destroy()
		log.Fatal(err)
		return
	}

	// wait for the process to finish.
	if _, err := process.Wait(); err != nil {
		log.Fatal(err)
	}

	// destroy the container.
	container.Destroy()
}

var (
	defaultMountFlags = unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV
	config            = &configs.Config{
		Capabilities: &configs.Capabilities{
			Bounding: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
			Effective: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
			Inheritable: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
			Permitted: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
			Ambient: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
		},
		Namespaces: configs.Namespaces([]configs.Namespace{
			{Type: configs.NEWNS},
			{Type: configs.NEWUTS},
			{Type: configs.NEWIPC},
			{Type: configs.NEWPID},
			{Type: configs.NEWUSER},
		}),
		Cgroups: &configs.Cgroup{
			Name:   "test-container",
			Parent: "system",
			Resources: &configs.Resources{
				MemorySwappiness: nil,
				AllowAllDevices:  nil,
				AllowedDevices:   configs.DefaultAllowedDevices,
			},
		},
		MaskPaths: []string{
			"/proc/kcore",
			"/sys/firmware",
		},
		ReadonlyPaths: []string{
			"/proc/sys", "/proc/sysrq-trigger", "/proc/irq", "/proc/bus",
		},
		Devices:  configs.DefaultAutoCreatedDevices,
		Hostname: "testing",
		Mounts: []*configs.Mount{
			{
				Source:      "proc",
				Destination: "/proc",
				Device:      "proc",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "tmpfs",
				Destination: "/dev",
				Device:      "tmpfs",
				Flags:       unix.MS_NOSUID | unix.MS_STRICTATIME,
				Data:        "mode=755",
			},
			{
				Source:      "devpts",
				Destination: "/dev/pts",
				Device:      "devpts",
				Flags:       unix.MS_NOSUID | unix.MS_NOEXEC,
				Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
			},
			{
				Device:      "tmpfs",
				Source:      "shm",
				Destination: "/dev/shm",
				Data:        "mode=1777,size=65536k",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "mqueue",
				Destination: "/dev/mqueue",
				Device:      "mqueue",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "sysfs",
				Destination: "/sys",
				Device:      "sysfs",
				Flags:       defaultMountFlags | unix.MS_RDONLY,
				Data:        "rbind",
			},
		},
		UidMappings: []configs.IDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        65536,
			},
		},
		GidMappings: []configs.IDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        65536,
			},
		},
		Rlimits: []configs.Rlimit{
			{
				Type: unix.RLIMIT_NOFILE,
				Hard: uint64(1025),
				Soft: uint64(1025),
			},
		},
	}
)
