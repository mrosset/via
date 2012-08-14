package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	return
	time.Sleep(time.Second / 2)
	r, e := http.Get("http://localhost:8080/")
	if e != nil {
		log.Fatal(e)
	}
	io.Copy(os.Stdout, r.Body)
}
