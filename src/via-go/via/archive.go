package via

import (
	"archive/tar"
	"compress/gzip"
	"debug/elf"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	Sharedlib  = "application/x-sharedlib"
	Executable = "application/x-executable"
)

func Package(name string, arch string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	defer func() {
		os.Chdir(wd)
	}()
	plan, err := FindPlan(name)
	if err != nil {
		return err
	}
	dir := filepath.Join(packages, plan.NameVersion())
	file := fmt.Sprintf("%s-%s.tar.gz", plan.NameVersion(), arch)
	plan.Tarball = file
	file = filepath.Join(repo, arch, file)
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	gz, err := gzip.NewWriterLevel(fd, gzip.BestCompression)
	if err != nil {
		return err
	}
	mani := new(Manifest)
	mani.Meta = (plan)
	vis := NewTarVisitor(gz, mani)
	walkFn := func(path string, info os.FileInfo, err error) error {
		vis.VisitFile(path, info)
		return nil
	}
	filepath.Walk(".", walkFn)
	err = WriteGzFile(mani, manifestName)
	if err != nil {
		return err
	}
	err = vis.tarFile(manifestName)
	if err != nil {
		return err
	}
	vis.tw.Close()
	gz.Close()
	fd.Close()

	return Sign(file)
}

type TarVisitor struct {
	tw        *tar.Writer
	man       *Manifest
	hardlinks map[uint64]string
}

func NewTarVisitor(w io.Writer, m *Manifest) *TarVisitor {
	return &TarVisitor{
		tar.NewWriter(w),
		m,
		make(map[uint64]string),
	}
}

func (tv TarVisitor) VisitDir(path string, f os.FileInfo) bool {
	if path == "." {
		return true
	}
	hdr := NewHeader(path, tv.hardlinks)
	tv.tw.WriteHeader(hdr)
	tv.man.AddEntry(path, EntryDir)
	return true
}

func (tv TarVisitor) VisitFile(path string, f os.FileInfo) {
	// TODO: remove this vpack does packaging
	if path == "DEPENDS" || path == "MANIFEST" || path == "manifest.json.gz" {
		return
	}
	var (
		deps []string
	)
	mime, err := fileMagic(path)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	switch mime {
	case Sharedlib:
		err = stripLib(path)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		deps, err = getDepends(path)
		if err != nil {
			fmt.Println("ERROR", err)
		}
	case Executable:
		err = stripBin(path)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		deps, err = getDepends(path)
		if err != nil {
			fmt.Println("ERROR", err)
		}
	}
	_ = deps
	err = tv.tarFile(path)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	tv.man.AddEntry(path, EntryFile)
}

func (tv TarVisitor) tarFile(path string) (err error) {
	hdr := NewHeader(path, tv.hardlinks)
	tv.tw.WriteHeader(hdr)
	if hdr.Typeflag == tar.TypeSymlink || hdr.Typeflag == tar.TypeLink {
		return nil
	}
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(tv.tw, fd)
	return err
}

func NewHeader(path string, hl map[uint64]string) (hdr *tar.Header) {
	hdr = new(tar.Header)
	fi, err := os.Lstat(path)
	if err != nil {
		fmt.Println(err)
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatal(path + " is not a Unix file")
	}
	hdr.Name = path
	hdr.Mode = int64(fi.Mode())
	hdr.Uid = int(stat.Uid)
	hdr.Gid = int(stat.Gid)
	hdr.AccessTime = time.Now()
	hdr.ModTime = time.Now()
	hdr.ChangeTime = time.Now()
	switch {
	case fi.IsDir():
		hdr.Typeflag = tar.TypeDir
		hdr.Name = hdr.Name + "/"
	case stat.Nlink > 1:
		if ref, ok := hl[stat.Ino]; ok {
			hdr.Typeflag = tar.TypeLink
			hdr.Linkname = ref
			break
		}
		hl[stat.Ino] = path
		hdr.Typeflag = tar.TypeReg
		hdr.Size = fi.Size()
	case fi.Mode() == os.ModeSymlink:
		link, err := os.Readlink(path)
		if err != nil {
			fmt.Println(err)
		}
		hdr.Typeflag = tar.TypeSymlink
		hdr.Linkname = link
	default:
		hdr.Typeflag = tar.TypeReg
		hdr.Size = fi.Size()
	}
	return hdr
}

type TarBallReader struct {
	fd *os.File
	gz *gzip.Reader
	tr *tar.Reader
}

func (this *TarBallReader) Close() {
	this.gz.Close()
	this.fd.Close()
}

func NewTarBallReader(path string) (tgzr *TarBallReader, err error) {
	fd, err := os.Open(path)
	if err != nil {
		return
	}
	gz, err := gzip.NewReader(fd)
	if err != nil {
		return
	}
	tr := tar.NewReader(gz)
	if err != nil {
		return
	}
	tgzr = &TarBallReader{fd, gz, tr}
	return
}

func Unpack(root string, file string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	defer func() {
		os.Chdir(wd)
	}()
	err = os.Chdir(root)
	if err != nil {
		return
	}
	tgr, err := NewTarBallReader(file)
	if err != nil {
		return err
	}
	defer tgr.Close()
	for {
		hdr, err := tgr.tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Name == manifestName {
			break
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			//fmt.Printf("\r%-40.40s -> D", hdr.Name)
			if fileExists(hdr.Name) {
				continue
			}
			err = os.Mkdir(hdr.Name, os.FileMode(hdr.Mode))
			if err != nil {
				fmt.Println(err)
				return err
			}
		case tar.TypeLink:
			if fileExists(hdr.Name) {
				err := os.Remove(hdr.Name)
				if err != nil {
					fmt.Println(err)
				}
			}
			err = os.Link(hdr.Linkname, hdr.Name)
			if err != nil {
				fmt.Println(err)
			}
		case tar.TypeSymlink:
			if fileExists(hdr.Name) {
				err = os.Remove(hdr.Name)
				if err != nil {
					fmt.Println(err)
				}
			}
			//fmt.Printf("\r%-40.40s -> %s", hdr.Name, hdr.Linkname)
			if err != nil {
				fmt.Println(err)
			}
			err = os.Symlink(hdr.Linkname, hdr.Name)
			if err != nil {
				fmt.Println(err)
			}
		case tar.TypeReg, tar.TypeRegA:
			//fmt.Printf("\r%-40.40s -> F", hdr.Name)
			f, err := os.OpenFile(hdr.Name, os.O_WRONLY|os.O_CREATE, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tgr.tr)
			f.Close()
			if err != nil {
				return err
			}
		}
	}
	return
}

func fileMagic(path string) (string, error) {
	output, err := exec.Command("file", "-b", "-i", path).Output()
	if err != nil {
		return "", err
	}
	mime := strings.Split(string(output), " ")
	return mime[0][:len(mime[0])-1], nil
}

func stripLib(path string) error {
	stripArg := "--strip-unneeded"
	output, err := exec.Command("strip", stripArg, path).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s %s", string(output), err)
	}
	return err
}

func stripBin(path string) error {
	stripArg := "--strip-all"
	output, err := exec.Command("strip", stripArg, path).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s %s", string(output), err)
	}
	return err
}

func getDepends(path string) (depends []string, err error) {
	f, err := elf.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	depends, err = f.ImportedLibraries()
	return depends, err
}
