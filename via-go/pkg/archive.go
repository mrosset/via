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
	wd, _ := os.Getwd()
	defer func() {
		os.Chdir(wd)
	}()
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
	err = manifest.Save(manifestName)
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
	if hdr.Typeflag == tar.TypeSymlink {
		return
	}
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
	fi, err := os.Lstat(path)
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
		link, err := os.Readlink(path)
		if err != nil {
			fmt.Println(err)
		}
		hdr.Name = path
		hdr.Typeflag = tar.TypeSymlink
		hdr.Linkname = link
	default:
		hdr.Typeflag = tar.TypeReg
		hdr.Size = fi.Size
	}
	return hdr
}

type TarBallReader struct {
	fd *os.File
	gz *gzip.Decompressor
	tr *tar.Reader
}

func (this *TarBallReader) Close() {
	fmt.Println("closeing tarball")
	this.gz.Close()
	this.fd.Close()
}

func NewTarBallReader(path string) (tgzr *TarBallReader, err os.Error) {
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

func Unpack(root string, file string) (err os.Error) {
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
		if err != nil {
			return
		}
		if hdr == nil {
			break
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			//fmt.Printf("%-40.40s -> D\n", hdr.Name)
			if fileExists(hdr.Name) {
				continue
			}
			err = os.Mkdir(hdr.Name, uint32(hdr.Mode))
			if err != nil {
				fmt.Println(err)
				return
			}
		case tar.TypeSymlink:
			if fileExists(hdr.Name) {
				err = os.Remove(hdr.Name)
				if err != nil {
					fmt.Println(err)
				}
			}
			//fmt.Printf("%-40.40s -> %s\n", hdr.Name, hdr.Linkname)
			if err != nil {
				fmt.Println(err)
			}
			err = os.Symlink(hdr.Linkname, hdr.Name)
			if err != nil {
				fmt.Println(err)
			}
		case tar.TypeReg, tar.TypeRegA:
			//fmt.Printf("%-40.40s -> F\n", hdr.Name)
			f, err := os.OpenFile(hdr.Name, os.O_WRONLY|os.O_CREATE, uint32(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(f, tgr.tr)
			f.Close()
			if err != nil {
				return
			}
		}
	}
	return
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fi.IsRegular() || fi.IsDirectory() {
		return true
	}
	return false
}
