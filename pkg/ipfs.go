package via

import (
	"context"
	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
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
	node, err := core.NewNode(context.TODO(), &core.BuildCfg{NilRepo: true})
	if err != nil {
		return "", err
	}
	fd, err := os.Open(path.String())
	if err != nil {
		return "", err
	}
	defer fd.Close()
	return coreunix.Add(node, fd)
}
