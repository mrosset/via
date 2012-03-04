package via

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"util/file"
)

var (
	ErrorTarHeader = errors.New("Unknown tar header")
)

func Peek(cr io.Reader) (dir string, err error) {
	tr := tar.NewReader(cr)
	hdr, err := tr.Next()
	if err != nil && err != io.EOF {
		return "", err
	}
	return path.Clean(hdr.Name), nil
}

func RmTar(file, dest string) (err error) {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	gz, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err != nil && err != io.EOF {
			return err
		}
		if hdr == nil {
			break
		}
		switch hdr.Typeflag {
		case tar.TypeReg:
			f := path.Join(dest, hdr.Name)
			if Verbose {
				info("Removeing", f)
			}
			err := os.Remove(f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Decompress Reader to destination directory
func Untar(cr io.Reader, dest string) (err error) {
	if !file.Exists(dest) {
		return fmt.Errorf("Directory %s does not exists.", dest)
	}
	tr := tar.NewReader(cr)
	for {
		hdr, err := tr.Next()
		if err != nil && err != io.EOF {
			return err
		}
		if hdr == nil {
			break
		}
		// Switch through header Typeflag and handle tar entry accordingly 
		switch hdr.Typeflag {
		// Handles Directories
		case tar.TypeDir:
			path := path.Join(dest, hdr.Name)
			if err := mkDir(path, hdr.Mode); err != nil {
				return err
			}
			continue
		case tar.TypeReg, tar.TypeRegA:
			path := path.Join(dest, hdr.Name)
			if err := writeFile(path, hdr, tr); err != nil {
				return err
			}
			continue
		default:
			fmt.Println(hdr.Name, "*** Unknown Header Type ***")
		}
	}
	return
}

func Tar(wr io.Writer, dir string) (err error) {
	tw := tar.NewWriter(wr)
	defer tw.Close()
	walkFn := func(path string, info os.FileInfo, err error) error {
		spath := strings.Replace(path, dir, "", -1)
		if spath == "" {
			return nil
		}
		spath = spath[1:]
		fi, err := os.Stat(path)
		hdr := fiToHeader(spath, fi)
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		switch {
		case hdr.Typeflag == tar.TypeDir:
		default:
			fd, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fd.Close()
			_, err = io.Copy(tw, fd)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return filepath.Walk(dir, walkFn)
}

func fiToHeader(name string, fi os.FileInfo) (hdr *tar.Header) {
	hdr = new(tar.Header)
	hdr.Name = name
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		log.Fatal(errors.New(fi.Name() + " is not a Unix file"))
	}
	hdr.Mode = int64(stat.Mode)
	hdr.Uid = int(stat.Uid)
	hdr.Gid = int(stat.Gid)
	hdr.AccessTime = time.Now()
	hdr.ModTime = time.Now()
	hdr.ChangeTime = time.Now()
	switch fi.IsDir() {
	case true:
		hdr.Typeflag = tar.TypeDir
	case false:
		hdr.Typeflag = tar.TypeReg
		hdr.Size = stat.Size
	}
	return hdr
}

// Make directory with permission
func mkDir(path string, mode int64) (err error) {
	if file.Exists(path) {
		return
	}
	info("mkdir", path)
	err = os.Mkdir(path, os.FileMode(mode))
	if err != nil {
		return err
	}
	return
}

// Write file from tar reader
func writeFile(path string, hdr *tar.Header, tr *tar.Reader) (err error) {
	if file.Exists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.FileMode(hdr.Mode))
	if err != nil {
		return err
	}
	if Verbose {
		info("Write", path)
	}
	//pb := console.NewProgressBarWriter(filepath.Base(path), hdr.Size, fd)
	_, err = io.Copy(fd, tr)
	fd.Close()
	if err != nil {
		return err
	}
	return
}
