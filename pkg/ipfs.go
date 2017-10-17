package via

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	IPFS_BIN = "ipfs"
)

func lookPathPanic(path string) string {
	arg0, err := exec.LookPath(path)
	if err != nil {
		panic(err)
	}
	return arg0
}

// Add 'path' to ipfs, returns content ID as a string
// TODO: use ipfs coreunix instead of this hackish exec
func Store(path Path) (string, error) {
	buf := new(bytes.Buffer)
	tee := io.MultiWriter(os.Stdout, buf)
	ipfs := exec.Cmd{
		Path: lookPathPanic("ipfs"),
		Args: []string{
			"ipfs", "add", "-rHn", "--local", path.String(),
		},
		Stdout: tee,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}

	if err := ipfs.Run(); err != nil {
		return "", err
	}
	scan := bufio.NewScanner(buf)
	last := p""
	// TODO: wrap this in a go func
	for scan.Scan() {
		last = scan.Text()
	}
	if scan.Err() != nil {
		return "", scan.Err()
	}
	s := strings.Split(last, " ")
	if len(s) != 3 {
		return "", fmt.Errorf("could not parse CID")
	}
	return s[1], scan.Err()
}
