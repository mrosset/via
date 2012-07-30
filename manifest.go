package via

import (
	"debug/elf"
	"errors"
	"fmt"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/file/magic"
	"github.com/str1ngs/util/json"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

func strip(p string) error {
	m, err := magic.GetFileMagic(p)
	if err != nil {
		return err
	}
	if m.Enum != magic.MagicElf {
		return nil
	}
	if verbose {
		fmt.Printf(lfmt, "strip", p)
	}
	cmd := exec.Command("strip", p, "--strip-all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateManifest(dir string, plan *Plan) (err error) {
	mfile := join(dir, "manifest.json.gz")
	files := []string{}
	if file.Exists(mfile) {
		err := os.Remove(mfile)
		if err != nil {
			elog.Println(err)
			return err
		}
	}
	var size int64
	walkFn := func(path string, info os.FileInfo, err error) error {
		if path == dir {
			return nil
		}

		spath := path[len(dir)+1:]
		stat, err := os.Lstat(path)
		if err != nil {
			elog.Println(err, path)
			return err
		}
		if stat.IsDir() {
			return nil
		}
		//strip(path)
		size = size + stat.Size()
		files = append(files, spath)
		return nil
	}
	err = filepath.Walk(dir, walkFn)
	if err != nil {
		return err
	}
	plan.Depends = Depends(plan.Name, dir, files)
	plan.Files = files
	plan.Date = time.Now()
	plan.Size = size
	plan.Save()
	return json.WriteGzJson(&plan, mfile)
}

func Depends(pname, base string, files []string) []string {
	deps := []string{}
	for _, j := range files {
		d := depends(join(base, j))
		for _, k := range d {
			o := owns(k)
			if o == "glibc" {
				continue
			}
			if o == "" {
				fmt.Println("warning", "can not resolve", k)
				continue
			}
			if contains(deps, o) || pname == o {
				continue
			}
			deps = append(deps, o)
		}
	}
	if len(deps) == 0 {
		return nil
	}
	return deps
}

func contains(sl []string, s string) bool {
	for _, j := range sl {
		if j == s {
			return true
		}
	}
	return false
}

func depends(file string) []string {
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

func owns(file string) string {
	e, err := filepath.Glob(join(config.Plans, "*.json"))
	if err != nil {
		elog.Println(err)
	}
	for _, j := range e {
		p, err := ReadPath(j)
		if err != nil {
			elog.Println(fmt.Errorf("%s %s", j, err))
			continue
		}
		for _, f := range p.Files {
			if filepath.Base(f) == file {
				return p.Name
			}
		}
	}
	return ""
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
