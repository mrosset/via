package via

import (
        "fmt"
        "github.com/mrosset/util/json"
        "sort"
)

// Repo provides repo path string. This is the path that binary
// tarballs are downloaded and built too
type Repo struct {
        Path
}

// NewRepo returns a new initialized Repo
func NewRepo(path string) Repo {
        return Repo{
                Path(path),
        }
}

// File returns the full path for repo.json file
func (r Repo) File(config *Config) Path {
        return config.Plans.Join("repo.json")
}

// FilesFile returns the full path for files.json
func (r Repo) FilesFile(config *Config) Path {
        return Path(config.Plans.Join("files.json"))
}

// RepoFiles provides plan files map hash
type RepoFiles map[string][]string

// Returns a sorted  key string slice
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
        if err := json.Read(&files, config.Repo.FilesFile(config).String()); err != nil {
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
        if err := json.Write(repo, config.Repo.File(config).String()); err != nil {
                return err
        }
        return json.Write(files, config.Repo.FilesFile(config).String())
}
