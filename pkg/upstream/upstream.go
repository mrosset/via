package upstream

import (
	"astuart.co/goq"
	"github.com/blang/semver"
	"net/http"
	"regexp"
	"strings"
	"text/scanner"
)

type nginx struct {
	Title string   `goquery:"h1"`
	Links []string `goquery:"pre a"`
}

type example struct {
	Title string   `goquery:"h1"`
	Files []string `goquery:"table.files tbody tr.js-navigation-item td.content,text"`
}

// func KernelMirror() error {
//	res, err := http.Get("https://ftp.gnu.org/gnu/bash/?C=M;O=D")
//	if err != nil {
//		return err
//	}
//	var n nginx

//	err = goq.NewDecoder(res.Body).Decode(&n)
//	if err != nil {
//		return err
//	}
//	f := n.Files
//	for _, v := range f {
//		fmt.Println(v)
//	}
//	return nil
// }

func Upstream() error {
	res, err := http.Get("http://mirrors.kernel.org/gnu/bash/")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var n nginx

	err = goq.NewDecoder(res.Body).Decode(&n)
	if err != nil {
		return err
	}

	for _, l := range n.Links {
		_, _ = semver.Parse(l)
	}
	return nil
}

func isInt(i rune) bool {
	switch i {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func ParseName(in string) string {
	var s scanner.Scanner
	s.Init(strings.NewReader(in))
	s.Mode = scanner.ScanChars
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == '-' && isInt(s.Peek()) {
			return (in[:s.Position.Column-1])
		}
	}
	return ""
}

func ParseVersion(in string) string {
	var (
		regs = []string{
			"[0-9]+.[0-9]+.[0-9]+.[0-9]+",
			"[0-9]+.[0-9]+.[0-9]+",
			"[0-9]+.[0-9]+",
		}
	)
	for _, r := range regs {
		v := regexp.MustCompile(r).FindString(in)
		if v != "" {
			return v
		}
	}
	return ""
}

func ParseNameVersion(file string) (string, string) {
	return ParseName(file), ParseVersion(file)
}
