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
}

func (t *Builder) RpcBuild(req Request, resp *Response) error {
	construct := NewConstruct(GetConfig(), &req.Plan)
	err := construct.BuildSteps()
	if err != nil {
		return err
	}
	return Install(construct.Config, req.Plan.Name)
}

func StartDaemon() error {
	rpc.Register(&Builder{})
	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		return err
	}
	if !IsInstalled("devel") {
		Install(GetConfig(), "devel")
	}
	defer os.Remove(SOCKET_FILE)
	go rpc.Accept(l)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	<-signals
	return nil
}
