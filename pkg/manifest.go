package via

import (
	"debug/elf"
	"errors"
	"fmt"
	"github.com/mrosset/util/json"
	"os"
	"os/exec"
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

// CreateManifest walks the plans PKGDIR and creates a gzipped manifest file.
func CreateManifest(config *Config, plan *Plan, dir string) (err error) {
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
		if contains(removes, "/"+spath) {
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
		size += fi.Size()
		files = append(files, spath)
		return nil
	}
	err = filepath.Walk(dir, walkFn)
	if err != nil {
		return err
	}
	plan.Files = files
	plan.AutoDepends, err = Depends(config, plan, dir)
	if err != nil {
		return err
	}
	plan.Date = time.Now()
	plan.Size = size
	if err := WritePlan(config, plan); err != nil {
		elog.Println(err)
		return err
	}
	// Remove the old Cid before writing the package tarball
	plan.Cid = ""
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

// Depends return the plan names, that each elf file depends on.
func Depends(config *Config, plan *Plan, dir string) ([]string, error) {
	var (
		depends = []string{}
	)
	rfiles, err := ReadRepoFiles(config)
	if err != nil {
		return nil, err
	}
	for _, f := range plan.Files {
		n := needs(join(dir, f))
		if len(n) == 0 {
			continue
		}
		for _, d := range n {
			// skip this file if this plan owns this file
			if filesContains(plan.Files, d) {
				continue
			}
			owner := rfiles.Owns(d)
			if !contains(depends, owner) && owner != "" {
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

// ReadManifest returns an installed Plan's manifest by name
func ReadManifest(config *Config, name string) (*Plan, error) {
	man := new(Plan)
	file := config.DB.Installed(config).Join(name, "manifest.json")
	if err := json.Read(man, file.String()); err != nil {
		return nil, err
	}
	return man, nil
}

// Readelf prints the dynamic libs and interop sections for the elf
// binary specified by path name
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
