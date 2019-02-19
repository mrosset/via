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

// NewPath returns a new Path with paths joined. If the new path does not exist panic
func NewPath(paths ...string) Path {
        np := Path(filepath.Join(paths...))
        if !np.Exists() {
                panic(os.ErrNotExist)
        }
        return np
}

// ToPath Converts to Path
func (p Path) ToPath() Path {
        return Path(p)
}

// Clone url to this Path
func (p Path) Clone(url string) error {
        return Clone(p, url)
}

// Base returns the Path's base
func (p Path) Base() Path {
        return Path(filepath.Base(string(p)))
}

// Exists return true if the Path path exists
func (p Path) Exists() bool {
        return file.Exists(p.String())
}

// Ensure that the Path directory path is created
func (p Path) Ensure() {
        if err := p.MkdirAll(); err != nil {
                panic(err)
        }
}

// Touch Path
func (p Path) Touch() error {
        return file.Touch(p.String())
}

// MkdirAll recursively makes Path directory
func (p Path) MkdirAll() error {
        if p.Exists() {
                return nil
        }
        return os.MkdirAll(p.String(), DirMask)
}

// Join path arguments with the Path as parent. This is like
// filepath.Join but with this Path type as the parent
func (p Path) Join(s ...string) Path {
        join := append([]string{string(p)}, s...)
        return Path(filepath.Join(join...))
}

// Expand returns the Path as a string that has had its environment
// variables expanded
func (p Path) Expand() Path {
        return Path(os.ExpandEnv(string(p)))
}

// Ext returns the Path's file extension
func (p Path) Ext() string {
        return filepath.Ext(string(p))
}

// RemoveAll remove this Path recursively
func (p Path) RemoveAll() error {
        return os.RemoveAll(string(p))
}

// ExpandToPath is like Expand but returns a Path type
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
