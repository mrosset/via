package via

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/repo"
	rconfig "github.com/ipfs/go-ipfs/repo/config"
	ds2 "github.com/ipfs/go-ipfs/thirdparty/datastore2"
	"github.com/mrosset/util/json"
	"os"
)

func AddR(path string) (string, error) {
	conf := rconfig.Config{}
	if err := json.Read(&conf, Path("$HOME/.ipfs/config").String()); err != nil {
		return "", err
	}
	json.WritePretty(conf, os.Stdout)
	fmt.Println(conf)
	r := &repo.Mock{
		C: conf,
		D: ds2.ThreadSafeCloserMapDatastore(),
	}
	node, err := core.NewNode(context.TODO(), &core.BuildCfg{
		Repo: r,
	})
	if err != nil {
		return "", err
	}
	return coreunix.AddR(node, path)
}
