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
	SOCKET_FILE = "/tmp/via/socket"
)

type Request struct {
	Plan Plan
}

type Response struct {
}

type Builder struct {
	config *Config
}

func (t *Builder) RpcBuild(req Request, resp *Response) error {
	Clean(req.Plan.Name)
	err := BuildSteps(t.config, &req.Plan)
	if err != nil {
		return err
	}
	in := NewInstaller(t.config, &req.Plan)
	return in.Install()
}

func StartDaemon(config *Config) error {
	rpc.Register(&Builder{config: config})
	if !file.Exists(filepath.Dir(SOCKET_FILE)) {
		os.Mkdir(filepath.Dir(SOCKET_FILE), 0700)
	}
	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		return err
	}
	if !IsInstalled(config, "devel") {
		p, err := NewPlan("devel")
		if err != nil {
			return err
		}
		batch := NewBatch(config)
		batch.Walk(p)
		batch.Install()
	}
	defer os.Remove(SOCKET_FILE)
	go rpc.Accept(l)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	<-signals
	return nil
}
