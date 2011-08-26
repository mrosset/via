package via

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func Package(name string, arch string) (err os.Error) {
	plan, err := FindPlan(name)
	if err != nil {
		return err
	}
	dir := filepath.Join(packages, plan.NameVersion())
	file := fmt.Sprintf("%s-%s.tar.gz", plan.NameVersion(), arch)
	file = filepath.Join(repo, arch, file)
	err = os.Chdir(dir)
	if err != nil {
		return
	}
	fd, err := os.Create(file)
	if err != nil {
		return
	}
	gz, err := gzip.NewWriterLevel(fd, gzip.BestCompression)
	if err != nil {
		return
	}
	manifest := NewManifest(plan)
	vis := NewTarVisitor(gz, manifest)
	filepath.Walk(".", vis, nil)
	err = manifest.Save("manifest.json.gz")
	if err != nil {
		return
	}
	err = tarFile("manifest.json.gz", vis.tw)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	vis.tw.Close()
	gz.Close()
	fd.Close()
	return
}

type TarVisitor struct {
	tw  *tar.Writer
	man *Manifest
}

func NewTarVisitor(w io.Writer, m *Manifest) *TarVisitor {
	return &TarVisitor{tar.NewWriter(w), m}
}

func (t TarVisitor) VisitDir(path string, f *os.FileInfo) bool {
	if path == "." {
		return true
	}
	hdr := NewHeader(path)
	t.tw.WriteHeader(hdr)
	t.man.AddEntry(path, EntryDir)
	return true
}

func (t TarVisitor) VisitFile(path string, f *os.FileInfo) {
	err := tarFile(path, t.tw)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	t.man.AddEntry(path, EntryFile)
}

func tarFile(path string, tw *tar.Writer) (err os.Error) {
	hdr := NewHeader(path)
	tw.WriteHeader(hdr)
	fd, err := os.Open(path)
	if err != nil {
		return
	}
	io.Copy(tw, fd)
	fd.Close()
	return
}

func NewHeader(path string) *tar.Header {
	hdr := new(tar.Header)
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	hdr.Name = path
	hdr.Mode = int64(fi.Mode)
	hdr.Uid = fi.Uid
	hdr.Gid = fi.Gid
	hdr.Atime = time.Seconds()
	hdr.Mtime = time.Seconds()
	hdr.Ctime = time.Seconds()
	switch {
	case fi.IsDirectory():
		hdr.Typeflag = tar.TypeDir
		hdr.Name = hdr.Name + "/"
	case fi.IsSymlink():
		fmt.Println("LINK", path)
	default:
		hdr.Typeflag = tar.TypeReg
		hdr.Size = fi.Size
	}
	return hdr
}
