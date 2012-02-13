package via

import (
	"fmt"
	"log"
	"os"
	"path"
	. "github.com/str1ngs/go-ansi/color"
	"path/filepath"
)

var (
	home     = "/home/strings/via"
	plans    = path.Join(home, "plans")
	repo     = path.Join(home, "repo")
	cache    = path.Join(home, "cache")
	packages = path.Join(cache, "packages")
	arch     = "x86_64"
	varDb    = "var/db/via"
)

func init() {
	log.SetPrefix(fmt.Sprintf("%s ", Blue("via:")))
	log.SetFlags(0)
}

const (
	PackExt = "tar.gz"
)

func GetRepo() string {
	return repo
}

func PkgFile(plan *Plan, arch string) string {
	return fmt.Sprintf("%s-%s-%s.%s", plan.Name, plan.Version, arch, PackExt)
}

func PkgAbsFile(plan *Plan, arch string) string {
	return fmt.Sprintf("%s/%s/%s", repo, arch, PkgFile(plan, arch))
}

func Install(root string, name string) os.Error {
	vd := filepath.Join(root, varDb, name)
	if !fileExists(vd) {
		err := os.MkdirAll(vd, 0755)
		if err != nil {
			return err
		}
	}
	plan, err := FindPlan(name)
	if err != nil {
		return err
	}
	tarball := PkgAbsFile(plan, arch)
	err = CheckSig(tarball)
	if err != nil {
		return err
	}
	mani, err := UnpackManifest(tarball)
	if err != nil {
		return err
	}
	log.Println("installing", name)
	err = Unpack(root, tarball)
	if err != nil {
		return err
	}
	err = WriteGzFile(mani, filepath.Join(vd, manifestName))
	return err
}

func List(root string, name string) os.Error {
	vd := filepath.Join(root, varDb, name)
	if !fileExists(vd) {
		return fmt.Errorf("%s name is not installed", name)
	}
	mani := new(Manifest)
	err := ReadGzFile(mani, filepath.Join(vd, manifestName))
	if err != nil {
		return err
	}
	for _, f := range mani.Files {
		fmt.Println(f.Path)
	}
	return nil
}

func Check(root string, name string) os.Error {
	vd := filepath.Join(root, varDb, name)
	if !fileExists(vd) {
		return fmt.Errorf("%s name is not installed", name)
	}
	mani := new(Manifest)
	err := ReadGzFile(mani, filepath.Join(vd, manifestName))
	if err != nil {
		return err
	}
	var errors = 0
	for _, f := range mani.Files {
		fpath := filepath.Join(root, f.Path)
		if !fileExists(fpath) {
			fmt.Println(fpath)
			errors++
			continue
		}
	}
	if errors != 0 {
		return fmt.Errorf("%d errors found", errors)
	}
	return nil
}

func Remove(root string, name string) os.Error {
	vd := filepath.Join(root, varDb, name)
	if !fileExists(vd) {
		return fmt.Errorf("%s name is not installed", name)
	}
	mani := new(Manifest)
	err := ReadGzFile(mani, filepath.Join(vd, manifestName))
	if err != nil {
		return err
	}
	log.Println("removeing", name)

	// remove files
	for _, f := range mani.Files {
		fpath := filepath.Join(root, f.Path)
		if f.EntryType == EntryFile && fileExists(fpath) {
			err := os.Remove(fpath)
			if err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(vd)
}

func OwnsFile(root, file string) (*Manifest, os.Error) {
	vd := filepath.Join(root, varDb)
	pkgs, err := filepath.Glob(filepath.Join(vd, "*"))
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		mani := new(Manifest)
		err := ReadGzFile(mani, filepath.Join(pkg, manifestName))
		if err != nil {
			return nil, err
		}
		for _, f := range mani.Files {
			if f.EntryType == EntryLink || f.EntryType == EntryFile {
				if filepath.Base(f.Path) == file {
					return mani, nil
				}
			}
		}
	}
	return nil, nil
}

func logf(v ...interface{}) {
	format := ""
	for i := 0; i < len(v); i++ {
		format = format + "%-10.10s "
	}
	log.Printf(format, v...)
}

func dirEmpty(path string) bool {
	files, _ := filepath.Glob(path + "/*")
	if len(files) > 0 {
		return false
	}
	return true
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fi.IsRegular() || fi.IsDirectory() || fi.IsSymlink() {
		return true
	}
	return false
}
