package via

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"util/file"
)

func WriteGzJson(v interface{}, file string) (err error) {
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	return WriteGzIo(v, fd)
}

func ReadGzJson(v interface{}, file string) (err error) {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	return ReadGzIo(v, fd)
}

func ReadGzIo(v interface{}, r io.Reader) (err error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()
	return json.NewDecoder(gz).Decode(v)
}

func WriteGzIo(v interface{}, w io.Writer) (err error) {
	gz := gzip.NewWriter(w)
	defer gz.Close()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, b, "", "\t")
	if err != nil {
		return err
	}
	_, err = io.Copy(gz, buf)
	return err
}

func WriteJson(v interface{}, path string) (err error) {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, b, "", "\t")
	if err != nil {
		return err
	}
	_, err = io.Copy(fd, buf)
	return err
}

func ReadJson(name string) (plan *Plan, err error) {
	plan = new(Plan)
	path := path.Join(config.Plans(), name+".json")
	if !file.Exists(path) {
		return nil, errors.New("Could not find plan " + name)
	}
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	err = json.NewDecoder(fd).Decode(plan)
	if err != nil {
		return nil, err
	}
	return plan, err
}
