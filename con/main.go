package main

import (
	"fmt"
	"github.com/mrosset/util/json"
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
	err := json.Read(containConfig, "config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(containConfig.Rootfs)
	if len(os.Args) > 1 && os.Args[1] == "init" {
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
	abs, err := filepath.Abs("rootfs")
	fmt.Println(abs)
	if err != nil {
		log.Fatal(err)
	}
	factory, err := libcontainer.New(abs, libcontainer.RootlessCgroupfs, libcontainer.InitArgs("/proc/self/exe", "init"))
	if err != nil {
		panic(err)
		return
	}
	container, err := factory.Create("container-id", containConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	process := &libcontainer.Process{
		Args:   []string{"/bin/bash"},
		Env:    []string{"PATH=/bin"},
		User:   "root",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = container.Run(process)
	if err != nil {
		container.Destroy()
		log.Fatal(err)
		return
	}

	// wait for the process to finish.
	_, err = process.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// destroy the container.
	container.Destroy()

}

var (
	defaultMountFlags = unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV
	containConfig     = &configs.Config{}
)
