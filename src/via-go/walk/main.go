package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	markFn := func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	}
	err := filepath.Walk("/", markFn)
	if err != nil {
		log.Fatal(err)
	}
}
