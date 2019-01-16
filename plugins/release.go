package main

import (
	"fmt"
	"github.com/ipfs/go-ipfs-api"
	"github.com/mrosset/util/file"
	"github.com/mrosset/via/pkg"
	"log"
	"os"
	"os/exec"
)

var (
	elog = log.New(os.Stderr, "", log.Lshortfile)
	lfmt = "%-20.20s %v\n"
)

type release struct {
	config *via.Config
}

func NewRelease(config *via.Config) release {
	return release{
		config: config,
	}
}

func (f *release) SetConfig(config *via.Config) {
	f.config = config
	f.config.Root = "/tmp/root"
	f.config.DB = "db"
}

func (f release) Execute() error {
	if f.config == nil {
		return fmt.Errorf("config is nil use SetConfig method first")
	}
	if f.config.Root != "/tmp/root" {
		return fmt.Errorf("config root is not sane")
	}
	defer os.RemoveAll("/tmp/root")
	fmt.Printf(lfmt, "executing", "release")
	if file.Exists(f.config.Repo) {
		fmt.Printf(lfmt, "cleaning", "repo")
		if err := os.RemoveAll(f.config.Repo); err != nil {
			return err
		}
	}
	plan, err := via.NewPlan("devel")
	if err != nil {
		return err
	}
	batch := via.NewBatch(f.config)
	batch.Walk(plan)

	errors := batch.Install()

	if len(errors) > 0 {
		log.Fatal(errors)
	}

	shell := shell.NewShell(f.config.IpfsApi)
	hash, err := shell.AddDir(f.config.Repo)
	if err != nil {
		return err
	}
	return execs("ipfs-cluster-ctl", "pin", "add", hash)
}

// Executes 'cmd' with 'args' useing os.Stdout and os.Stderr
func execs(cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Stderr = os.Stderr
	e.Stdout = os.Stdout
	return e.Run()
}

var Release release
