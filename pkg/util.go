package via

// This file contains private functions reused throughout the Api

import (
        "fmt"
        "github.com/mrosset/util/file"
        "github.com/mrosset/util/json"
        "os"
        "path/filepath"
)

// function aliases
var (
        join   = filepath.Join
        base   = filepath.Base
        exists = file.Exists
)

func baseDir(p string) string {
        return base(filepath.Dir(p))
}

// Returns true if a string slice contains a string
func contains(sl []string, s string) bool {
        for _, j := range sl {
                if expand(j) == expand(s) {
                        return true
                }
        }
        return false
}

// Checks if a dir path is empty
func isEmpty(p string) bool {
        e, err := filepath.Glob(join(p, "*"))
        if err != nil {
                elog.Println(err)
                return false
        }
        return len(e) == 0
}

// Walks a path and returns a slice of dir/files
func walkPath(p string) ([]string, error) {
        s := []string{}
        fn := func(path string, fi os.FileInfo, err error) error {
                s = append(s, path)
                return nil
        }
        err := filepath.Walk(p, fn)
        if err != nil {
                return nil, err
        }
        return s, nil
}

func cpDir(s, d string) error {
        pdir := filepath.Dir(s)
        fn := func(p string, fi os.FileInfo, err error) error {
                spath := p[len(pdir)+1:]
                dpath := join(d, spath)
                if file.Exists(dpath) {
                        return nil
                }
                if fi.IsDir() {
                        return os.Mkdir(dpath, fi.Mode())
                }
                fd, err := os.OpenFile(dpath, os.O_CREATE|os.O_WRONLY, fi.Mode())
                if err != nil {
                        elog.Println(err)
                        return err
                }
                defer fd.Close()
                return file.Copy(fd, p)
        }
        return filepath.Walk(s, fn)
}

func printSlice(s []string) {
        for _, j := range s {
                fmt.Println(j)
        }
}

// PlanFilePath returns the full path of the Plan's json file
func PlanFilePath(config *Config, plan *Plan) Path {
        return config.Plans.Join(plan.Group, plan.Name+".json")
}

// WritePlan writes the serialized go struct to it's json file. The
// json file is pretty formatted so to keep consistency
func WritePlan(config *Config, plan *Plan) error {
        file := PlanFilePath(config, plan).String()
        return json.Write(PlanJSON(*plan), file)
}
