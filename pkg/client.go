package via

import (
	"net"
	"net/rpc"
)

func Connect() (*rpc.Client, error) {
	l, err := net.Dial("unix", SOCKET_FILE)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(l), nil
}
