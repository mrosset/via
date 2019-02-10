package via

import (
	"fmt"
	"github.com/mrosset/util/json"
	"sort"
)

type RepoFiles map[string][]string

// Returns a sorted slice key strings
func (rf RepoFiles) keys() []string {
	var (
		keys = []string{}
	)
	for k := range rf {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Returns the first alphabetical plan Name of plan that contains file
func (rf RepoFiles) Owns(file string) string {
	for _, key := range rf.keys() {
		if filesContains(rf[key], file) {
			return key
		}
	}
	fmt.Println("warning: can not resolve", file)
	return ""
}

// Like Owns but returns a slice of plan names instead of the first
// occurrence. The returned slice is sorted alphabetically
func (rf RepoFiles) Owners(file string) []string {
	s := []string{}
	for _, key := range rf.keys() {
		if filesContains(rf[key], file) {
			s = append(s, key)
		}
	}
	return s
}

func ReadRepoFiles(config *Config) (RepoFiles, error) {
	files := RepoFiles{}
	if err := json.Read(&files, join(config.Plans, "files.json")); err != nil {
		return nil, err
	}
	return files, nil
}

func RepoCreate(config *Config) error {
	var (
		repo  = []string{}
		files = map[string][]string{}
		rfile = join(config.Plans, "repo.json")
		ffile = join(config.Plans, "files.json")
	)
	e, err := PlanFiles()
	if err != nil {
		return err
	}
	for _, j := range e {
		p, err := ReadPath(config, j)
		if err != nil {
			return err
		}
		repo = append(repo, join(p.Group, p.Name+".json"))
		files[p.Name] = p.Files
	}
	err = json.Write(repo, rfile)
	if err != nil {
		return err
	}
	return json.Write(files, ffile)
}
