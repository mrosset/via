// +build ignore

package via

import (
	"context"
	"github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/commands/files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/repo"
	rconfig "github.com/ipfs/go-ipfs/repo/config"
	ds2 "github.com/ipfs/go-ipfs/thirdparty/datastore2"
	"github.com/mrosset/util/json"
	"os"
	"path/filepath"
)

func AddR(path string) (string, error) {
	conf := rconfig.Config{}
	if err := json.Read(&conf, Path("$HOME/.ipfs/config").String()); err != nil {
		return "", err
	}
	_ = &repo.Mock{
		C: conf,
		// C: rconfig.Config{
		//	Identity: rconfig.Identity{
		//		PeerID: "Qmfoo", // required by offline node
		//	},
		// },
		D: ds2.ThreadSafeCloserMapDatastore(),
	}

	node, err := core.NewNode(context.TODO(), &core.BuildCfg{
		NilRepo: true,
	})

	exch := node.Exchange
	bserv := blockservice.New(node.Blockstore, exch)
	dserv := merkledag.NewDAGService(bserv)

	if err != nil {
		return "", err
	}
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
