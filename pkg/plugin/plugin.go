package plugin

import (
	"github.com/mrosset/via/pkg"
	"os"
	"os/exec"
	"path/filepath"
)

type Plugin interface {
	SetConfig(*via.Config)
	Execute() error
}

// Executes 'cmd' with 'args' useing os.Stdout and os.Stderr
func execs(cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Stderr = os.Stderr
	e.Stdout = os.Stdout
	return e.Run()
}

func build(out string, in string) error {
	return execs("go", "build", "-buildmode=plugin", "-o", out, in)
}

func Build(config *via.Config) error {
	dir := filepath.Join(config.Repo, "../plugins")
	glob := filepath.Join(dir, "*.go")
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, in := range files {
		out := in[:len(in)-3] + ".so"
		if err := build(out, in); err != nil {
			return err
		}

	}
	return nil
}
