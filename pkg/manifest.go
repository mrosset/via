package via

import (
	"debug/elf"
	"errors"
	"fmt"
	"github.com/str1ngs/util/json"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"time"
)

func strip(p string) error {
	ef, err := elf.Open(p)
	// if elf.Open fails then its not a elf file skip it.
	if err != nil {
		return nil
	}
	defer ef.Close()
	if verbose {
		fmt.Printf(lfmt, os.Getenv("STRIP"), base(p))
	}
	cmd := exec.Command("strip", p)
	return cmd.Run()
}

// Walk the package directory and make a file list.
// The resulting file list and plan data, is saved
// to manifest.json which then gets tar/gzipped into
// the package file.
//
// CreateManifest also perform strip and removale of
// blacklisted files.
func CreateManifest(dir string, plan *Plan) (err error) {
	var (
		size  int64
		mfile = join(dir, "manifest.json.gz")
		files = []string{}
	)
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if path == dir {
			return nil
		}
		// FIXME: Do removes in Package
		spath := path[len(dir)+1:]
		removes := append(config.Remove, plan.Remove...)
		// If the file is in config.Remove or plan.Removes delete it
		if contains(removes, spath) {
			// TODO: expand path
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
			fmt.Printf(lfmt, "removing", spath)
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if fi.IsDir() {
			return nil
		}
		//FIXME: stripping breaks gcc packaging
		//if err := strip(path); err != nil {
		//return err
		//}
		// We need to restat after strip.
		si, err := os.Lstat(path)
		if err != nil {
			return err
		}
		size += si.Size()
		files = append(files, spath)
		return nil
	}
	err = filepath.Walk(dir, walkFn)
	if err != nil {
		return err
	}
	plan.Files = files
	plan.Depends, err = Depends(dir, plan)
	if err != nil {
		return err
	}
	plan.Date = time.Now()
	plan.Size = size
	plan.Save()
	return json.WriteGz(&plan, mfile)
}

func filesContains(files []string, file string) bool {
	for _, f := range files {
		if base(f) == file {
			return true
		}
	}
	return false
}

func Depends(dir string, plan *Plan) ([]string, error) {
	depends := []string{}
	rfiles, err := ReadRepoFiles()
	if err != nil {
		return nil, err
	}
	for _, f := range plan.Files {
		n := needs(join(dir, f))
		if len(n) == 0 {
			continue
		}
		for _, d := range n {
			if filesContains(plan.Files, d) {
				// skip this file if this plan owns this file
				continue
			}
			owner := rfiles.Owns(d)
			if !contains(depends, owner) {
				depends = append(depends, owner)
			}
		}
	}
	sort.Strings(depends)
	return depends, nil
}

func needs(file string) []string {
	f, err := elf.Open(file)
	if err != nil {
		return nil
	}
	im, err := f.ImportedLibraries()
	if err != nil {
		return nil
	}
	return im
}

func ReadManifest(name string) (man *Plan, err error) {
	man = new(Plan)
	err = json.Read(man, path.Join(config.DB.Installed(), name, "manifest.json"))
	if err != nil {
		return
	}
	return
}

func Readelf(p string) error {
	f, err := elf.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()
	im, err := f.ImportedLibraries()
	if err != nil {
		return err
	}
	fmt.Printf(lfmt, "libs", im)
	sec := f.Section(".interp")
	d, err := sec.Data()
	if err != nil {
		return err
	}
	fmt.Printf(lfmt, "intr", string(d))
	ds := f.SectionByType(elf.SHT_DYNAMIC)
	d, err = ds.Data()
	if err != nil {
		return err
	}
	str, err := stringTable(f, ds.Link)
	if err != nil {
		return err
	}
	for len(d) > 0 {
		// TODO: add byteorder for ELFCLASS32
		tag := elf.DynTag(f.ByteOrder.Uint64(d[0:8]))
		val := uint64(f.ByteOrder.Uint64(d[8:16]))
		d = d[16:]
		if tag == elf.DT_RPATH {
			s, ok := getString(str, int(val))
			if ok {
				fmt.Printf(lfmt, "rpath", s)
			}
		}
	}
	return nil
}

// FIXME: These 2 functions are taken from GOROOT/src/pkg/elf.
// add license or request they be exported?
// getString extracts a string from an ELF string table.
func getString(section []byte, start int) (string, bool) {
	if start < 0 || start >= len(section) {
		return "", false
	}

	for end := start; end < len(section); end++ {
		if section[end] == 0 {
			return string(section[start:end]), true
		}
	}
	return "", false
}

// stringTable reads and returns the string table given by the
// specified link value.
func stringTable(f *elf.File, link uint32) ([]byte, error) {
	if link <= 0 || link >= uint32(len(f.Sections)) {
		return nil, errors.New("section has invalid string table link")
	}
	return f.Sections[link].Data()
}
