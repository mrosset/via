package via

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mikesun/go-ipfs-api"
	files "github.com/mikesun/go-multipart-files"
	"github.com/mrosset/util/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	IPFS_BIN = "ipfs"
)

var (
	IPFS_API_FILE = Path("$HOME/.ipfs/api")
	IPFS_API      = "http://localhost:5001/api/v0"
)

func lookPathPanic(path string) string {
	arg0, err := exec.LookPath(path)
	if err != nil {
		panic(err)
	}
	return arg0
}

type IpfsStat struct {
	Path    string
	Mode    os.FileMode
	ModTime time.Time
}

func readApi() string {
	api, err := ioutil.ReadFile(IPFS_API_FILE.String())
	if err != nil {
		elog.Fatalf("could not get api from %s", IPFS_API_FILE)
	}
	return string(api)
}

func NewLocalShell() *shell.Shell {
	return shell.NewShell(IPFS_API)
}

type IpfsResponse struct {
	Version string
}

func IpfsVersion() (string, error) {
	vr := &IpfsResponse{}
	if err := json.Get(vr, IPFS_API+"/version"); err != nil {
		return "", err
	}
	return vr.Version, nil
}

func AddR(path Path) (string, error) {
	fi, err := os.Stat(path.String())
	if err != nil {
		return "", err
	}
	sf, err := files.NewSerialFile("", path.String(), fi)
	if err != nil {
		return "", err
	}
	slf := files.NewSliceFile("", path.String(), []files.File{sf})
	mpr := files.NewMultiFileReader(slf, true)
	len, err := slf.Size()
	if err != nil {
		return "", err
	}
	req, err := NewMultiPartRequest(IPFS_API+"/add?recursive=true", len, mpr)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
		return "", err
	}
	io.Copy(os.Stderr, res.Body)
	return "", nil
}

func NewMultiPartRequest(url string, len int64, r *files.MultiFileReader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+r.Boundary())
	req.Header.Set("Content-Disposition", "form-data: name=\"files\"")
	return req, nil
}
func OAddR(path Path) (string, error) {
	s := NewLocalShell()
	return s.AddDir(path.String())
}

func Add(path Path) (string, error) {
	s := NewLocalShell()
	fd, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	return s.Add(fd)
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
		files = append(files, IpfsStat{Path: p, Mode: info.Mode(), ModTime: info.ModTime()})
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
		fpath := path.JoinS(f.Path)
		if err := os.Chmod(string(fpath), f.Mode); err != nil {
			return err
		}
		os.Chtimes(fpath.String(), time.Now(), f.ModTime)
	}
	return nil
}

// Add 'path' to ipfs, returns content ID as a string
// TODO: use ipfs coreunix instead of this hackish exec
func IpfsAdd(path Path, hashOnly bool) (string, error) {
	var (
		buf   = new(bytes.Buffer)
		tee   = io.MultiWriter(os.Stdout, buf)
		flags = "-rH"
		scan  = bufio.NewScanner(buf)
		last  = ""
	)
	if hashOnly {
		flags += "n"
	}
	ipfs := exec.Cmd{
		Path: lookPathPanic("ipfs"),
		Args: []string{
			"ipfs", "add", flags, "--local", path.String(),
		},
		Stdout: tee,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}
	if path.IsDir() {
		if err := MakeStat(path); err != nil {
			return "", err
		}
	}
	if err := ipfs.Run(); err != nil {
		return "", err
	}
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
