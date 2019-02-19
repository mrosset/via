package via

import (
	"net"
	"net/rpc"
)

// Connect dials the rpc daemon
func Connect() (*rpc.Client, error) {
	l, err := net.Dial("unix", SocketFile)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(l), nil
}
