package via

import (
	"fmt"
	"os"
)

const (
	manifestName = "manifest.json.gz"
)

const (
	EntryDir = iota
	EntryFile
	EntryLink
)

type FileEntry struct {
	Path      string
	EntryType int
}

type Manifest struct {
	Meta  *Plan
	Files []*FileEntry
}

func (t *Manifest) AddEntry(file string, eType int) {
	t.Files = append(t.Files, &FileEntry{file, eType})
}

func UnpackManifest(file string) (mani *Manifest, err os.Error) {
	tbr, err := NewTarBallReader(file)
	mani = new(Manifest)
	if err != nil {
		return nil, err
	}
	defer tbr.Close()
	for {
		hdr, err := tbr.tr.Next()
		if err == os.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Name == manifestName {
			err = ReadGzIo(mani, tbr.tr)
			return mani, err
		}
	}
	return nil, fmt.Errorf("Could not find %s in %s", manifestName, file)
}
