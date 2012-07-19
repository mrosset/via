package via

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrorTarHeader = errors.New("Unknown tar header")
)

// TODO: rewrite this hackfest
// Decompress Reader to destination directory
func Untar(r io.Reader, dest string) (man *Manifest, err error) {
	if !file.Exists(dest) {
		return nil, fmt.Errorf("%s does not exist.", dest)
	}
	man = new(Manifest)
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			elog.Println(err)
			return nil, err
		}
		path := join(dest, hdr.Name)
		// Switch through header Typeflag and handle tar entry accordingly 
		switch hdr.Typeflag {
		// Handles Directories
		case tar.TypeDir:
			if err := mkDir(path, hdr.Mode); err != nil {
				return nil, err
			}
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
					return nil, err
				}
				continue
			}
			if err != nil && err != io.EOF {
				return nil, err
			}
			// Write long file data to disk
			if err := writeFile(fpath, hdr, tr); err != nil {
				return nil, err
			}
		case tar.TypeSymlink:
			err := os.Symlink(hdr.Linkname, path)
			if err != nil {
				elog.Fatal(err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if hdr.Name == "manifest.json.gz" {
				err := json.ReadGzIo(man, tr)
				if err != nil {
					return nil, err
				}
				continue
			}
			dir := filepath.Dir(path)
			if !file.Exists(dir) {
				fmt.Println(dir)
				elog.Println("FIXME: (hdr permission) tar has no top directory.")
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					elog.Println(err)
					return nil, err
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
	return
}

// TODO: rewrite this hackfest
func Tarball(wr io.Writer, plan *Plan) (err error) {
	dir := join(cache.Pkgs(), plan.NameVersion())
	err = CreateManifest(dir, plan)
	if err != nil {
		return err
	}
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
		err = tw.WriteHeader(hdr)
		if err != nil {
			elog.Println(err)
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir, tar.TypeSymlink:
		case tar.TypeReg:
			fd, err := os.Open(path)
			if err != nil {
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
