package via

import (
	"archive/tar"
	"bytes"
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
	"util/json"
)

var (
	ErrorTarHeader = errors.New("Unknown tar header")
)

type TarGzReader struct {
	fd *os.File
	gz *gzip.Reader
	Tr *tar.Reader
}

func (tgzr *TarGzReader) Close() {
	tgzr.gz.Close()
	tgzr.fd.Close()
}

func NewTarGzReader(pfile string) (tgzr *TarGzReader, err error) {
	fd, err := os.Open(pfile)
	if err != nil {
		return nil, err
	}
	gz, err := gzip.NewReader(fd)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	return &TarGzReader{fd, gz, tr}, nil
}

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
func Untar(r io.Reader, dest string) (err error) {
	if !file.Exists(dest) {
		return fmt.Errorf("Directory %s does not exists.", dest)
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err != nil && err != io.EOF {
			return err
		}
		if hdr == nil {
			break
		}
		// Switch through header Typeflag and handle tar entry accordingly 
		switch {
		// Handles Directories
		case hdr.Typeflag == tar.TypeDir:
			path := path.Join(dest, hdr.Name)
			if err := mkDir(path, hdr.Mode); err != nil {
				return err
			}
			continue
		case string(hdr.Typeflag) == "L":
			lfile := new(bytes.Buffer)
			// Get longlink path from tar file data
			lfile.ReadFrom(tr)
			fpath := path.Join(dest, lfile.String())
			// Read next iteration for file data
			hdr, err := tr.Next()
			if hdr.Typeflag == tar.TypeDir {
				err := mkDir(fpath, hdr.Mode)
				if err != nil {
					return err
				}
				continue
			}
			if err != nil && err != io.EOF {
				return err
			}
			// Write long file data to disk
			if err := writeFile(fpath, hdr, tr); err != nil {
				return err
			}
		case hdr.Typeflag == tar.TypeReg, hdr.Typeflag == tar.TypeRegA:
			path := path.Join(dest, hdr.Name)
			if err := writeFile(path, hdr, tr); err != nil {
				return err
			}
			continue
		default:
			fmt.Println(hdr.Name, "*** Unknown Header Type ***")
		}
		continue
	}
	return
}

func Package(wr io.Writer, plan *Plan) (err error) {
	man := &Manifest{Plan: plan}
	dir := config.GetPackageDir(plan.NameVersion())
	tw := tar.NewWriter(wr)
	defer tw.Close()
	walkFn := func(path string, info os.FileInfo, err error) error {
		spath := strings.Replace(path, dir, "", -1)
		if spath == "" {
			return nil
		}
		fe := &File{Path: spath}
		spath = spath[1:]
		fi, err := os.Stat(path)
		hdr := fiToHeader(spath, fi)
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		switch {
		case hdr.Typeflag == tar.TypeDir:
			fe.Type = TypeDir

		default:
			fe.Type = TypeFile
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
		man.Files = append(man.Files, fe)
		return nil
	}
	err = filepath.Walk(dir, walkFn)
	if err != nil {
		return err
	}
	json.Write(man, "manifest.json")
	file.Cat("manifest.json")
	return nil
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
	if _, err = io.Copy(fd, tr); err != nil {
		return err
	}
	fd.Close()
	if err != nil {
		return err
	}
	return
}
