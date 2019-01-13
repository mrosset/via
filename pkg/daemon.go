package via

import (
	"net"
	"net/rpc"
	"os"
	"os/signal"
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
	return Install(t.config, req.Plan.Name)
}

func StartDaemon(config *Config) error {
	rpc.Register(&Builder{config: config})
	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		return err
	}
	if !IsInstalled(config, "devel") {
		Install(config, "devel")
	}
	defer os.Remove(SOCKET_FILE)
	go rpc.Accept(l)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	<-signals
	return nil
}
