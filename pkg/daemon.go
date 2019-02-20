package via

import (
	"github.com/mrosset/util/file"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	//SocketFile is full path to socket file
	SocketFile = "/tmp/via/socket"
)

// Request provide RPC request type
type Request struct {
	Plan Plan
}

// Response is RPC response type
type Response struct {
}

// DaemonBuilder provides RPC server type
//
// FIXME: this is not complete and not as important since we are using namespaces.
// could be useful at a later point
type DaemonBuilder struct {
	config *Config
}

// RPCBuild calls the RPC to build a plan
// func (t *DaemonBuilder) RPCBuild(req Request, _ *Response) error {
//      ctx := NewPlanContext(t.config, &req.Plan)
//      Clean(ctx)
//      if err := BuildSteps(ctx); err != nil {
//              return err
//      }
//      return NewInstaller(t.config, &req.Plan).Install()
// }

// StartDaemon starts the RPC daemon
func StartDaemon(config *Config) error {
	rpc.Register(&DaemonBuilder{config: config})
	if !file.Exists(filepath.Dir(SocketFile)) {
		os.Mkdir(filepath.Dir(SocketFile), 0700)
	}
	l, err := net.Listen("unix", SocketFile)
	if err != nil {
		return err
	}
	if !IsInstalled(config, "devel") {
		p, err := NewPlan(config, "devel")
		if err != nil {
			return err
		}
		batch := NewBatch(config)
		batch.Walk(p)
		batch.Install()
	}
	defer os.Remove(SocketFile)
	go rpc.Accept(l)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	<-signals
	return nil
}
