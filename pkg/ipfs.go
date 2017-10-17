package via

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mrosset/util/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	IPFS_BIN = "ipfs"
)

func lookPathPanic(path string) string {
	arg0, err := exec.LookPath(path)
	if err != nil {
		panic(err)
	}
	return arg0
}

type IpfsStat struct {
	Path    Path
	Mode    os.FileMode
	ModTime time.Time
}

// Walk 'path' and creates a stat.json file with each files permissions
func MakeStat(path Path) error {
	var (
		files = []IpfsStat{}
		sfile = path.JoinS("stat.json")
	)
	fn := func(p string, info os.FileInfo, err error) error {
		// if this is root the directory or the stat file do nothing
		if p == path.String() || p == sfile.String() {
			return nil
		}
		p = strings.Replace(p, path.String()+"/", "", 1)
		files = append(files, IpfsStat{Path: Path(p), Mode: info.Mode(), ModTime: info.ModTime()})
		return nil
	}
	filepath.Walk(path.String(), fn)
	return json.Write(files, sfile.String())
}

// Sets each files Mode in 'path' to mode contained in the paths stat.json file
func SetStat(path Path) error {
	var (
		files = []IpfsStat{}
		sfile = path.JoinS("stat.json")
	)
	if !sfile.Exists() {
		return fmt.Errorf("%s: does not have a stat.json file", path)
	}
	if err := json.Read(&files, sfile.String()); err != nil {
		return err
	}
	for _, f := range files {
		fpath := path.Join(f.Path)
		if err := os.Chmod(fpath.String(), f.Mode); err != nil {
			return err
		}
		os.Chtimes(fpath.String(), time.Now(), f.ModTime)
	}
	return nil
}

// Add 'path' to ipfs, returns content ID as a string
// TODO: use ipfs coreunix instead of this hackish exec
func IpfsAdd(path Path) (string, error) {
	buf := new(bytes.Buffer)
	tee := io.MultiWriter(os.Stdout, buf)
	ipfs := exec.Cmd{
		Path: lookPathPanic("ipfs"),
		Args: []string{
			"ipfs", "add", "-rH", "--local", path.String(),
		},
		Stdout: tee,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}
	MakeStat(path)
	if err := ipfs.Run(); err != nil {
		return "", err
	}
	scan := bufio.NewScanner(buf)
	last := ""
	// TODO: wrap this in a go func
	for scan.Scan() {
		last = scan.Text()
	}
	if scan.Err() != nil {
		return "", scan.Err()
	}
	s := strings.Split(last, " ")
	if len(s) != 3 {
		return "", fmt.Errorf("could not parse CID")
	}
	return s[1], scan.Err()
}

func IpfsGet(dir Path, cid string) error {
	ipfs := exec.Cmd{
		Path: lookPathPanic("ipfs"),
		Args: []string{
			"ipfs", "get", cid,
		},
		Dir:    dir.String(),
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}
	if err := ipfs.Run(); err != nil {
		return err
	}
	return SetStat(dir.JoinS(cid))
}
