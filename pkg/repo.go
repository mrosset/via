package via

import (
	"fmt"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"os"
	"path/filepath"
	"sort"
)

// Repo provides repo path string. This is the path that binary
// tarballs are downloaded and built too
type Repo string

// String provides stringer interface
func (r Repo) String() string {
	return os.ExpandEnv(string(r))
}

// Exists return true if the Repo path exists
func (r Repo) Exists() bool {
	return file.Exists(r.String())
}

// Ensure that the directory is created
func (r Repo) Ensure() error {
	if r.Exists() {
		return nil
	}
	return os.MkdirAll(r.String(), 0755)
}

// Expand returns the Repo path as a string that has been its
// environmental variables expanded.
func (r Repo) Expand() string {
	return os.ExpandEnv(string(r))
}

// NewRepo returns a new Repo who's parent is joined with dir
func NewRepo(parent, dir string) Repo {
	return Repo(filepath.Join(parent, dir))
}

// RepoFiles provides plan files map hash
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

// Owns returns the first alphabetical plan Name of plan that contains file
func (rf RepoFiles) Owns(file string) string {
	for _, key := range rf.keys() {
		if filesContains(rf[key], file) {
			return key
		}
	}
	fmt.Println("warning: can not resolve", file)
	return ""
}

// Owners like owns but returns a slice of plan names instead of the first
// occurrence. The returned slice is sorted alphabetically
func (rf RepoFiles) Owners(file string) []string {
	owners := []string{}
	for _, key := range rf.keys() {
		if filesContains(rf[key], file) {
			owners = append(owners, key)
		}
	}
	return owners
}

// ReadRepoFiles reads files.json and returns a RepoFiles map hash
func ReadRepoFiles(config *Config) (RepoFiles, error) {
	files := RepoFiles{}
	if err := json.Read(&files, join(config.Plans, "files.json")); err != nil {
		return nil, err
	}
	return files, nil
}

// RepoCreate reads each plan's files creating a repo.json file that
// contains all plan's and groups. And also creating a files.json that
// contains a hash map of each plans files
//
// FIXME: this is pretty expensive and probably won't scale well. Also
// repo.json and files.json should probably not be kept in version control.
func RepoCreate(config *Config) error {
	var (
		repo  = []string{}
		files = map[string][]string{}
		rfile = join(config.Plans, "repo.json")
		ffile = join(config.Plans, "files.json")
	)
	e, err := PlanFiles(config)
	if err != nil {
		return err
	}
	for _, j := range e {
		p, err := ReadPath(j)
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
