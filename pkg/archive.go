package via

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrorTarHeader = errors.New("Unknown tar header")
)

func GNUUntar(dest string, file string) error {
	tar := exec.Command("tar", "-xf", file)
	tar.Dir = dest
	tar.Stdout = os.Stdout
	tar.Stderr = os.Stdout
	return tar.Run()
}

// TODO: rewrite this hackfest
// Decompress Reader to destination directory
func Untar(dest string, r io.Reader) error {
	if !file.Exists(dest) {
		return fmt.Errorf("%s does not exist.", dest)
	}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			elog.Println(err)
			return err
		}
		if hdr.Name == "manifest.json.gz" {
			continue
		}
		//fmt.Printf("%c %s\n", hdr.Typeflag, hdr.Name)
		path := join(dest, hdr.Name)
		// Switch through header Typeflag and handle tar entry accordingly
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := mkDir(path, hdr.Mode); err != nil {
				return err
			}
			// Long File
		case 'L':
			lfile := new(bytes.Buffer)
			// Get longlink path from tar file data
			lfile.ReadFrom(tr)
			fpath := join(dest, lfile.String())
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
		case tar.TypeSymlink:
			err := os.Symlink(hdr.Linkname, path)
			if err != nil {
				elog.Fatal(err)
			}
		case tar.TypeReg, tar.TypeRegA:
			dir := filepath.Dir(path)
			if !file.Exists(dir) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					elog.Println(err)
					return err
				}
			}
			if err := writeFile(path, hdr, tr); err != nil {
				elog.Println(err)
			}
			continue
		default:
			fmt.Println(hdr.Name, "*** Unknown Header Type ***")
		}
		continue
	}
	return nil
}

// Walk directory and tars files to io.Writer
func archive(wr io.Writer, dir string) error {
	tw := tar.NewWriter(wr)
	defer tw.Close()
	walkFn := func(path string, info os.FileInfo, err error) error {
		spath := strings.Replace(path, dir, "", -1)
		if spath == "" {
			return nil
		}
		spath = spath[1:]
		fi, err := os.Lstat(path)
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			elog.Println(err)
			return err
		}
		if hdr.Typeflag == tar.TypeSymlink {
			ln, err := os.Readlink(path)
			if err != nil {
				return err
			}
			hdr.Linkname = ln
		}
		hdr.Name = spath
		// TODO: check tar specs for actual length
		// If path is greater then 100 bytes we need to handle as LongLink
		if len(hdr.Name) >= 100 {
			hdr.Typeflag = tar.TypeGNULongName
			hdr.Size = int64(len(hdr.Name))
		}
		//fmt.Printf("%c %s\n", hdr.Typeflag, hdr.Name)
		err = tw.WriteHeader(hdr)
		if err != nil {
			elog.Println(err)
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir, tar.TypeSymlink:
		case tar.TypeGNULongName: // Handle long file paths
			// Write path as tar data.
			_, err := tw.Write([]byte(hdr.Name))
			if err != nil {
				elog.Println(err)
				return err
			}
			// Treat the long link as a file, flush so we can write the real data.
			tw.Flush()
			// Write a header so the writer knows the size of the data.
			hdr.Size = fi.Size()
			hdr.Typeflag = tar.TypeReg
			tw.WriteHeader(hdr)
			// Finally write the file to tar
			fd, err := os.Open(path)
			if err != nil {
				elog.Println(err)
				return err
			}
			defer fd.Close()
			_, err = io.Copy(tw, fd)
			if err != nil {
				elog.Println(err)
				return err
			}
		case tar.TypeReg:
			fd, err := os.Open(path)
			if err != nil {
				elog.Println(err)
				return err
			}
			defer fd.Close()
			_, err = io.Copy(tw, fd)
			if err != nil {
				elog.Println(err)
				return err
			}
		default:
			err = fmt.Errorf("%s: unhandled tar header type")
			elog.Println(err)
			return err
		}
		return nil
	}
	return filepath.Walk(dir, walkFn)

}

// TODO: rewrite this hackfest
func Tarball(wr io.Writer, plan *Plan) (err error) {
	dir := join(cache.Pkgs(), plan.NameVersion())
	err = CreateManifest(dir, plan)
	if err != nil {
		return err
	}
	return archive(wr, dir)
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

	//fmt.Printf(lfmt, "file", join("+ ", path))
	//pb := console.NewProgressBarWriter(path, hdr.Size, fd)
	if _, err = io.Copy(fd, tr); err != nil {
		return err
	}
	//pb.Close()
	fd.Close()
	if err != nil {
		return err
	}
	return
}

func TarGzReader(p string) (*tar.Reader, error) {
	fd, err := os.Open(p)
	if err != nil {
		elog.Println(err)
		return nil, err
	}
	gz, err := gzip.NewReader(fd)
	if err != nil {
		elog.Println(err)
		return nil, err
	}
	return tar.NewReader(gz), nil
}
func ReadPackManifest(p string) (*Plan, error) {
	man := new(Plan)
	tr, err := TarGzReader(p)
	if err != nil {
		elog.Println(err)
		return nil, err
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			elog.Println(err)
			return nil, err
		}
		if hdr.Name == "manifest.json.gz" {
			err := json.ReadGzIo(man, tr)
			if err != nil {
				elog.Println(err)
				return nil, err
			}
			return man, err
		}
	}
	return nil, fmt.Errorf("%s: could not find manifest.", p)
}
