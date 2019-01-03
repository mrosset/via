package via

import (
	"github.com/ipfs/go-ipfs-api"
	"os"
)

func IpfsAdd(path Path) (string, error) {
	s := shell.NewLocalShell()
	fd, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return s.Add(fd)
}

func HashOnly(path Path) (string, error) {
	s := shell.NewLocalShell()
	fd, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return s.Add(fd, shell.OnlyHash(true))
}
