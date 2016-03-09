package via

import (
	"os"
	"os/exec"
)

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
	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
