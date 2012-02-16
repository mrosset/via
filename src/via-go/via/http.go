package via

import (
	"bytes"
	"exp/html"
	"fmt"
	gurl "github.com/str1ngs/gurl/pkg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	client   *http.Client
	filesUrl = "https://via-test.googlecode.com/files"
	listUrl  = "https://code.google.com/p/via-test/downloads/list"
	netrc    = make(map[string]string)
)

func init() {
	if client == nil {
		client = new(http.Client)
		var err error
		netrc, err = getNetRc()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func upExists(file string) bool {
	url := fmt.Sprintf("%s/%s", filesUrl, file)
	res, _ := client.Head(url)
	if res.StatusCode == 200 {
		return true
	}
	return false
}

// Upload file to google code
func upload(file string) (err error) {
	if upExists(filepath.Base(file)) {
		fmt.Println("WARNING", file, "exists on server")
		//return nil
	}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	if err := w.WriteField("label", "label here"); err != nil {
		return err
	}
	if err := w.WriteField("summary", "summary here"); err != nil {
		return err
	}
	// Create file field writer
	fw, err := w.CreateFormFile("upload", filepath.Base(file))
	if err != nil {
		return err
	}
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd) //Write file part
	if err != nil {
		return err
	}
	// Important if you do not close the multipart writer you will not have a 
	// terminating boundry 
	w.Close()
	req, err := http.NewRequest("POST", filesUrl, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.SetBasicAuth(netrc["login"], netrc["password"])
	res, err := client.Do(req)
	return checkResponse(res, err)
}

func GetDownloadList() (list []string, err error) {
	res, err := client.Get(listUrl)
	if checkResponse(res, err) != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(buf)
	if err != nil {
		return nil, err
	}
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Data == "td" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "vt" {
					for _, c := range n.Child {
						if c.Data == "a" {
							list = append(list, filepath.Base(c.Attr[0].Val))
						}
					}
				}
			}
		}
		for _, c := range n.Child {
			fn(c)
		}
	}
	fn(doc)
	return list, err
}

func DownloadSrc(url string) (err error) {
	return download(url, cache)
}

func DownloadSig(url string) (err error) {
	return download(url+".sig", cache)
}

func download(url string, dest string) (err error) {
	gurl := new(gurl.Client)
	file := filepath.Base(url)
	_, err = os.Stat(filepath.Join(cache, file))
	if err == nil {
		log.Println(file + " exists skipping")
		return nil
	}
	return gurl.Download(cache, url)
}

func getNetRc() (map[string]string, error) {
	nr := make(map[string]string)
	home := os.Getenv("HOME")
	fd, err := os.Open(filepath.Join(home, ".netrc"))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(fd); err != nil {
		return nil, err
	}
	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = line[:len(line)-1]
		kv := strings.Split(string(line), " ")
		nr[kv[0]] = kv[1]
	}
	return nr, nil
}
