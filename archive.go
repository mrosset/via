package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"util"
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

// Decompress Reader to destination directory
func Untar(cr io.Reader, dest string) (err error) {
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
			fmt.Println(hdr.Name, " *** Unknown Header Type **")
		}
	}
	return
}

// Make directory with permission
func mkDir(path string, mode int64) (err error) {
	if util.FileExists(path) {
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
	if util.FileExists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.FileMode(hdr.Mode))
	if err != nil {
		return err
	}
	_, err = io.Copy(f, tr)
	f.Close()
	if err != nil {
		return err
	}
	return
}
