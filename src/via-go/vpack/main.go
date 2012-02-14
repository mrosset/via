package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"via-go/via"
)

var (
	arch = flag.String("arch", "x86_64", "which cpu architect to use")
	root = flag.String("root", "/", "specify the root directory")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	switch args[0] {
	case "install":
		args = shift(args)
		err := install(args)
		if err != nil {
			log.Fatal(err)
		}
	case "pack":
		args = shift(args)
		err := pack(args)
		if err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
	}
}

func pack(targets []string) (err error) {
	for _, target := range targets {
		log.Printf("packing %s ", target)
		err = via.Package(target, *arch)
		if err != nil {
			return err
		}
	}
	return
}

func install(targets []string) (err error) {
	if !fileExists(*root) {
		err = os.MkdirAll(*root, 0755)
		if err != nil {
			return err
		}
	}
	for _, target := range targets {
		plan, err := via.FindPlan(target)
		if err != nil {
			return err
		}
		log.Printf("installing %s", plan.NameVersion())
		file := fmt.Sprintf("%s-%s.tar.gz", plan.NameVersion(), *arch)
		file = filepath.Join(via.GetRepo(), *arch, file)
		err = via.Unpack(*root, file)
		if err != nil {
			return err
		}
	}
	return
}

func shift(a []string) []string {
	return append(a[:0], a[0+1:]...)
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.IsDir() || fi.IsDir() {
		return true
	}
	return false
}
