package via

import (
	"fmt"
	cclient "github.com/gocircuit/circuit/client"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	FMT_CIRCUIT_ADDRESS = "circuit://%s/%s/%s"
	CIRCUIT_IP          = "192.168.0.22:1122"
)

// Search for a running circuit in /tmp and returns its formatted circuit address
func CircuitAddress() (string, error) {
	cs, _ := filepath.Glob("/tmp/circuit-*")
	if len(cs) != 1 {
		return "", fmt.Errorf("found %d running circuits should have one", len(cs))
	}
	sp := strings.Split(cs[0], "-")
	if len(sp) != 3 {
		return ",", fmt.Errorf("Malformed circuit lock path")
	}
	return fmt.Sprintf(FMT_CIRCUIT_ADDRESS, CIRCUIT_IP, sp[2][1:], sp[1]), nil
}

func pickServer(c *cclient.Client) cclient.Anchor {
	for _, r := range c.View() {
		return r
	}
	panic(0)
}

func CircuitBuild(name string) error {
	_, err := CircuitAddress()
	if err != nil {
		return err
	}
	c := cclient.DialDiscover("228.8.8.8:8822", nil)
	//c := cclient.Dial(a, nil)
	cmd := cclient.Cmd{
		Path:  "/home/strings/gocode/bin/via",
		Args:  []string{"build", "-c", name},
		Scrub: true,
	}
	t := pickServer(c).Walk([]string{"via", "build"})
	p, err := t.MakeProc(cmd)
	if err != nil {
		return err
	}
	p.Stdin().Close()
	go io.Copy(os.Stdout, p.Stdout())
	go io.Copy(os.Stderr, p.Stderr())
	_, err = p.Wait()
	return err
}
