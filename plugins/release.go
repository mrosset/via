package main

import (
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

// func NewRelease(config *via.Config) release {
//	return release{
//		config: config,
//	}
// }

func (f *release) SetConfig(config *via.Config) {
}

func (f release) Execute() error {
	return nil
}

// Executes 'cmd' with 'args' useing os.Stdout and os.Stderr
func execs(cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Stderr = os.Stderr
	e.Stdout = os.Stdout
	return e.Run()
}

// Release exports release type
var Release release
