package via

import (
	// "context"
	"fmt"
	"github.com/ipfs/go-ipfs-api"
	// "github.com/ipfs/go-ipfs-files"
	// "github.com/ipfs/go-ipfs/core"
	// "github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DockerEnv is the docker environment file
	DockerEnv = "/.dockerenv"
	// DockerAPI is the internal docker address used for ipfs API calls
	DockerAPI = "172.17.0.1:5001"
)

// IpfsStat provides type that contains file stat information
type IpfsStat struct {
	Path    string
	Mode    os.FileMode
	ModTime time.Time
}

func isDocker() bool {
	return file.Exists(DockerEnv)
}

func whichAPI(config *Config) string {
	if isDocker() {
		return DockerAPI
	}
	return config.IpfsAPI
}

// IpfsAdd add a file to ipfs and returns the ipfs multihash
func IpfsAdd(config *Config, path string) (string, error) {
	s := shell.NewShell(whichAPI(config))
	fd, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return s.Add(fd)
}

// HashOnly returns the ipfs multihash for a file at path
func HashOnly(config *Config, path string) (string, error) {
	s := shell.NewShell(whichAPI(config))
	fd, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return s.Add(fd, shell.OnlyHash(true))
}

// func CoreHashOnly(path Path) (string, error) {
//	node, err := core.NewNode(context.TODO(), &core.BuildCfg{Online: false}) // NilRepo: true})
//	if err != nil {

//		return "", err
//	}
//	fd, err := os.Open(path.String())
//	if err != nil {
//		return "", err
//	}
//	defer fd.Close()
//	adder, err := coreunix.NewAdder(context.TODO(), node.Pinning, node.Blockstore, node.DAG)
//	if err != nil {

//		return "", err
//	}
//	file := files.NewReaderFile(fd)
//	if err != nil {
//		return "", err
//	}
//	fn, err := adder.AddAllAndPin(file)
//	if err != nil {
//		return "", err
//	}
//	return fn.Cid().String(), nil
// }

// MakeStat walks path and creates a stat.json file with each files permissions
func MakeStat(path string) error {
	var (
		files = []IpfsStat{}
		sfile = filepath.Join(path, "stat.json")
	)
	fn := func(p string, info os.FileInfo, err error) error {
		// if this is root the directory or the stat file do nothing
		if p == path || p == sfile {
			return nil
		}
		p = strings.Replace(p, path+"/", "", 1)
		files = append(files, IpfsStat{Path: p, Mode: info.Mode(), ModTime: info.ModTime()})
		return nil
	}
	filepath.Walk(path, fn)
	return json.Write(files, sfile)
}

// SetStat each files Mode in path to mode contained in the paths stat.json file
func SetStat(path string) error {
	var (
		files = []IpfsStat{}
		sfile = filepath.Join(path, "stat.json")
	)
	if !file.Exists(sfile) {
		return fmt.Errorf("%s: does not have a stat.json file", path)
	}
	if err := json.Read(&files, sfile); err != nil {
		return err
	}
	for _, f := range files {
		fpath := filepath.Join(path, f.Path)
		if err := os.Chmod(string(fpath), f.Mode); err != nil {
			return err
		}
		os.Chtimes(fpath, time.Now(), f.ModTime)
	}
	return nil
}
