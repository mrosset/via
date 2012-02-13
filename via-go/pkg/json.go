package via

import (
	"bytes"
	"compress/gzip"
	"io"
	"json"
	"os"
)

func WriteGzFile(v interface{}, file string) (err os.Error) {
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	return WriteGzIo(v, fd)
}

func ReadGzFile(v interface{}, file string) (err os.Error) {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	return ReadGzIo(v, fd)
}

func ReadGzIo(v interface{}, r io.Reader) (err os.Error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()
	return json.NewDecoder(gz).Decode(v)
}

func WriteGzIo(v interface{}, w io.Writer) (err os.Error) {
	gz, err := gzip.NewWriter(w)
	if err != nil {
		return err
	}
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
