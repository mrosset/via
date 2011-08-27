package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"via"
)

var (
	arch = flag.String("arch", "x86_64", "which cpu architect to use")
	root = flag.String("root", "./tmp", "specify the root directory")
)

func init() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	args := flag.Args()
	switch args[0] {
	case "install":
		args = shift(args)
		err := install(args)
		if err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
	}
}

func pack(name, arch string) (err os.Error) {
	err = via.Package(name, arch)
	return
}

func install(targets []string) (err os.Error) {
	for _, target := range targets {
		plan, err := via.FindPlan(target)
		if err != nil {
			return err
		}
		fmt.Printf("install %-10.10s", plan.NameVersion())
		file := fmt.Sprintf("%s-%s.tar.gz", plan.NameVersion(), *arch)
		file = filepath.Join(via.GetRepo(), *arch, file)
		err = via.Unpack(*root, file)
		if err != nil {
			return err
		}
		fmt.Printf("done\n")
	}
	return
}

func shift(a []string) []string {
	return append(a[:0], a[0+1:]...)
}
