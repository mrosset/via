package via

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Pack Plan

type Plan struct {
	Name    string "name"
	Version string "version"
	Source  string "source"
	Tarball string "tarball"
}

func (this Plan) NameVersion() string {
	return fmt.Sprintf("%s-%s", this.Name, this.Version)
}

func (this Plan) Print() {
	pp := func(f, v string) {
		fmt.Printf("%-10.10s = %s\n", f, v)
	}
	pp("Name", this.Name)
	pp("Version", this.Version)
	pp("Source", this.Source)
	pp("Tarball", this.Tarball)
}

func ParsePlan(path string) (plan *Plan, err error) {
	var (
		kmap = make(map[string]string)
		keys = []string{"name", "version", "source"}
	)
	fd, err := os.Open(path)
	if err != nil {
		return
	}
	defer fd.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(fd)
	if err != nil {
		return
	}
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = line[:len(line)-1]
		for _, k := range keys {
			if strings.HasPrefix(line, k) {
				v := strings.Split(line, "=")[1]
				kmap[k] = v[1 : len(v)-1]
			}
		}
	}
	plan = &Plan{
		kmap["name"],
		kmap["version"],
		kmap["source"],
		"",
	}
	return
}

func FindPlan(name string) (plan *Plan, err error) {
	glob := fmt.Sprintf("%s/*/%s/plan", plans, name)
	files, err := filepath.Glob(glob)
	if err != nil {
		return
	}
	if len(files) != 1 {
		return nil, fmt.Errorf("expected 1 plan got %v", len(files))
	}
	plan, err = ParsePlan(files[0])
	return
}

func ListPlans() (err error) {
	glob := fmt.Sprintf("%s/*/*/plan", plans)
	files, err := filepath.Glob(glob)
	if err != nil {
		return
	}
	for _, i := range files {
		plan, err := ParsePlan(i)
		if err != nil {
			return fmt.Errorf("%s %s", i, err.Error())
		}
		fmt.Printf("%-20.20s %s\n", plan.Name, plan.Version)
	}
	fmt.Printf("found %v plans", len(files))
	return
}
