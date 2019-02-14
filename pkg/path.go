package via

import (
	"encoding/json"
	"github.com/mrosset/util/file"
	"os"
	"path/filepath"
)

const (
	// DirMask is the default mask for new directories
	DirMask = 0755
)

// Path provides type for working with directory paths
type Path string

// String provides stringer interface
func (p Path) String() string {
	return string(p)
}

// Exists return true if the Path path exists
func (p Path) Exists() bool {
	return file.Exists(p.String())
}

// Ensure that the Path directory path is created
func (p Path) Ensure() error {
	if p.Exists() {
		return nil
	}
	return os.MkdirAll(p.String(), DirMask)
}

// Join path arguments with the Path as parent. This is like
// filepath.Join but with this Path type as the parent
func (p Path) Join(s ...string) string {
	return filepath.Join(
		append([]string{string(p)}, s...)...,
	)
}

// Expand returns the Path as a string that has been its
// environmental variables expanded.
func (p Path) Expand() string {
	return os.ExpandEnv(string(p))
}

// ExpandToPath like
func (p Path) ExpandToPath() Path {
	return Path(p.Expand())
}

// MarshalJSON provide marshal interface
func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}

// UnmarshalJSON provide unmarshal interface
func (p *Path) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*p = Path(s)
	return nil
}
