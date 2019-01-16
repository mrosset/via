package via

import (
	"github.com/ipfs/go-ipfs-api"
	"github.com/mrosset/util/file"
	"os"
)

const (
	DOCKERENV = "/.dockerenv"
	DOCKERAPI = "172.17.0.1:5001"
)

func isDocker() bool {
	return file.Exists(DOCKERENV)
}

func whichApi() string {
	if isDocker() {
		return DOCKERAPI
	}
	return config.IpfsApi
}

func IpfsAdd(path Path) (string, error) {
	s := shell.NewShell(whichApi())
	fd, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return s.Add(fd)
}

func HashOnly(path Path) (string, error) {
	s := shell.NewShell(whichApi())
	fd, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return s.Add(fd, shell.OnlyHash(true))
}
