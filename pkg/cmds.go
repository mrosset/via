package via

import (
	"os"
	"os/exec"
)

func clone(dest, url string) error {
	cmd := exec.Command("git", "clone", url)
	cmd.Dir = dest
	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func wget(dest, url string) {
	cmd := exec.Command("wget", url)
	cmd.Dir = dest
	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func unzip(dest, file string) {
	cmd := exec.Command("unzip", file)
	cmd.Dir = dest
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
