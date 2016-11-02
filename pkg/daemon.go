package via

import (
	"github.com/valyala/gorpc"
	"os"
)

const (
	SOCKET = "/tmp/via.socket"
)

// Returns a rpc dispatcher with functions
func NewDispatcher() *gorpc.Dispatcher {
	d := gorpc.NewDispatcher()
	d.AddFunc("ping", ping)
	d.AddFunc("build", rpcBuild)
	return d
}

// Returns a rpc client and its dispatcher
func NewRpcClient() (*gorpc.Client, *gorpc.DispatcherClient) {
	var (
		d  = NewDispatcher()
		c  = gorpc.NewUnixClient(SOCKET)
		dc = d.NewFuncClient(c)
	)
	return c, dc
}

// Exportable RPC function for testing
func ping() bool {
	return true
}

func rpcBuild(p string) error {
	elog.Println(os.Getenv("PATH"))
	err := Install("devel")
	if err != nil {
		return err
	}
	plan, err := NewPlan(p)
	if err != nil {
		return err
	}
	return BuildSteps(plan)
}

func Listen() error {
	os.Remove(SOCKET)
	s := gorpc.NewUnixServer(SOCKET, NewDispatcher().NewHandlerFunc())
	s.Concurrency = 1
	return s.Serve()
}
