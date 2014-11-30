package github

import (
	"bitbucket.org/strings/via/pkg"
	"fmt"
	xfile "github.com/str1ngs/util/file"
	"github.com/str1ngs/util/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

//"https://uploads.github.com/repos/$owner/$repo/releases/$id/assets?name=$name"

var (
	config  = via.GetConfig()
	client  = new(http.Client)
	elog    = log.New(os.Stderr, "", log.Lshortfile)
	release = GetRelease()
)

const (
	token       = "09e1c6162e09441b36ea33bfb596cabbac6bd0a4"
	api_release = "https://api.github.com/repos/str1ngs/tarball/releases/743297"
	api_upload  = "https://uploads.github.com/repos/str1ngs/tarball/releases/743297/assets?name=%s"
)

type Release struct {
	Url    string
	Assets Assets
}

type Asset struct {
	Url  string
	Name string
}

type Assets []Asset

func Push(name string) error {
	plan, err := via.FindPlan(name)
	if err != nil {
		elog.Println(err)
		return err
	}
	url := fmt.Sprintf(api_upload, plan.PackageFile())
	file := path.Join(config.Repo, plan.PackageFile())
	if !xfile.Exists(file) {
		return nil
	}
	/*
		curl -H "Authorization: token $token" \
			-H "Content-Type: application/gzip" \
			--data-binary @$path\
			"https://uploads.github.com/repos/$owner/$repo/releases/$id/assets?name=$name"
	*/
	curl := &exec.Cmd{
		Path: "/usr/bin/curl",
		Args: []string{
			"-#",
			"-H", fmt.Sprintf("Authorization: token %s", token),
			"-H", "Content-Type: application/gzip",
			"--data-binary", fmt.Sprintf("@%s", file),
			url}}
	curl.Stdout = os.Stdout
	curl.Stderr = os.Stderr
	return curl.Run()
}

func PushOld(name string) error {
	plan, err := via.FindPlan(name)
	if err != nil {
		elog.Println(err)
		return err
	}
	url := fmt.Sprintf(api_upload, plan.PackageFile())
	file := path.Join(config.Repo, plan.PackageFile())
	elog.Println(file)
	fd, err := os.Open(file)
	if err != nil {
		elog.Println(err)
		return err
	}
	stat, _ := os.Stat(file)
	length := fmt.Sprintf("%d", stat.Size)
	elog.Println(length)
	req, err := http.NewRequest("POST", url, fd)
	req.Header.Add("Content-Type", "application/gzip")
	req.Header.Add("Content-Length", length)
	resp, err := client.Do(req)
	if err != nil {
		elog.Println(err)
		return err
	}
	io.Copy(os.Stdout, resp.Body)
	fmt.Println()
	elog.Println(plan.PackageFile())
	return nil
}

func PushAll() error {
	plans, err := via.GetPlans()
	if err != nil {
		return err
	}

	for _, i := range plans {
		if release.Assets.Contains(i.PackageFile()) {
			continue
		}
		if err := Push(i.Name); err != nil {
			elog.Println(err)
			return err
		}
	}
	return nil
}

func (a Assets) Contains(name string) bool {
	for _, i := range a {
		if i.Name == name {
			return true
		}
	}
	return false
}

func GetRelease() *Release {
	r := new(Release)
	err := json.Get(r, api_release)
	if err != nil {
		log.Fatal(err)
	}
	return r
}
