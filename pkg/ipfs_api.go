// +build ignore

package via

import (
	"context"
	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/commands/files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"os"
	"path/filepath"
)

func OAddR(path string) (string, error) {
	r, err := fsrepo.Open(Path("$HOME/.ipfs").String())
	if err != nil {
		panic(err)
		return "", err
	}
	node, err := core.NewNode(context.TODO(), &core.BuildCfg{
		Repo: r,
	})
	if err != nil {
		return "", err
	}
	bserv := blockservice.New(node.Blockstore, node.Exchange)
	dserv := merkledag.NewDAGService(bserv)
	return BaseAddR(node, path, dserv)
}

func BaseAddR(n *core.IpfsNode, root string, ds merkledag.DAGService) (key string, err error) {
	n.Blockstore.PinLock().Unlock()

	stat, err := os.Lstat(root)
	if err != nil {
		return "", err
	}

	f, err := files.NewSerialFile(filepath.Dir(root), root, false, stat)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fileAdder, err := coreunix.NewAdder(n.Context(), n.Pinning, n.Blockstore, ds)
	if err != nil {
		return "", err
	}
	err = fileAdder.AddFile(f)
	if err != nil {
		return "", err
	}

	nd, err := fileAdder.Finalize()
	if err != nil {
		return "", err
	}

	return nd.String(), nil
}
