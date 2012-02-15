package main

import (
	"flag"
	"fmt"
	. "github.com/str1ngs/go-ansi/color"
	"log"
	"via-go/via"
)

var (
	arch = flag.String("arch", "x86_64", "which cpu architect to use")
	root = flag.String("root", "/", "specify the root directory")
)

func init() {
	log.SetPrefix(fmt.Sprintf("%s ", Blue("via:")))
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	//args := flag.Args()
	//args = shift(args)
	url := "http://mirrors.kernel.org/gnu/bash/bash-4.2.tar.gz"
	log.Println("building", "bash")
	via.DownloadSrc(url)
}

func shift(a []string) []string {
	return append(a[:0], a[0+1:]...)
}
