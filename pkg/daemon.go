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
	Clean(req.Plan.Name)
	err := BuildSteps(&req.Plan)
	if err != nil {
		return err
	}
	return Install(req.Plan.Name)
}

func StartDaemon() error {
	rpc.Register(&Builder{})
	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		return err
	}
	if !IsInstalled("devel") {
		Install("devel")
	}
	defer os.Remove(SOCKET_FILE)
	go rpc.Accept(l)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	<-signals
	return nil
}