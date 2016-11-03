package via

import (
	"fmt"
	"github.com/valyala/gorpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

const (
	SOCKET = "/tmp/via.socket"
)

type server struct{}

func (s *server) Build(ctx context.Context, in *BuildRequest) (*BuildReply, error) {
	plan, err := NewPlan(in.Name)
	if err != nil {
		return nil, err
	}
	if in.Clean {
		err = Clean(plan.Name)
		if err != nil {
			return nil, err
		}
	}
	return &BuildReply{Message: "Finished building " + plan.NameVersion()}, BuildSteps(plan)
}
func Listen() error {
	os.Remove(SOCKET)
	lis, err := net.Listen("unix", SOCKET)
	if err != nil {
		return (err)
	}
	s := grpc.NewServer()
	RegisterBuilderServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
func unixDialer(addr string, timeout time.Duration) (net.Conn, error) {
	sock, err := net.DialTimeout("unix", SOCKET, timeout)
	return sock, err
}

func ClientRequestBuild(name string, clean bool) error {
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithDialer(unixDialer))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewBuilderClient(conn)

	r, err := c.Build(context.Background(), &BuildRequest{Name: name, Clean: clean})
	if err != nil {
		return fmt.Errorf("could not greet: %v", err)
	}
	elog.Println(r.Message)
	return nil
}

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
	err := Clean(p)
	if err != nil {
		return err
	}
	plan, err := NewPlan(p)
	if err != nil {
		return err
	}
	elog.Println("building", plan.NameVersion())
	err = BuildSteps(plan)
	if err != nil {
		elog.Println(err)
		return err
	}
	elog.Printf("done building %s", plan.NameVersion())
	return Install(p)
}

func OListen() error {
	os.Remove(SOCKET)
	s := gorpc.NewUnixServer(SOCKET, NewDispatcher().NewHandlerFunc())
	s.Concurrency = 1
	return s.Serve()
}
