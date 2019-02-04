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
	Files []string `goquery:"pre a"`
}

type example struct {
	Title string   `goquery:"h1"`
	Files []string `goquery:"table.files tbody tr.js-navigation-item td.content,text"`
}

func GnuUpstreamLatest(name, url string, current semver.Version) (string, error) {
	var (
		res    *http.Response
		err    error
		latest = "0.0.0"
	)
	if res, err = http.Get(url); err != nil {
		return "", err
	}
	defer res.Body.Close()

	n := new(nginx)

	if err = goq.NewDecoder(res.Body).Decode(&n); err != nil {
		return "", err
	}

	for _, l := range n.Files {
		n, v := ParseNameVersion(l)
		if n != name || v == "" {
			continue
		}
		sv, err := semver.ParseTolerant(v)
		if err != nil {
			continue
		}
		lv, err := semver.ParseTolerant(latest)
		if err != nil {
			continue
		}
		if sv.GT(current) && sv.GT(lv) {
			latest = v
		}

	}
	return latest, nil
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
