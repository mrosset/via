package upstream

import (
	"astuart.co/goq"
	"github.com/blang/semver"
	"net/http"
	"regexp"
	"strings"
	"text/scanner"
)

var (
	client = new(http.Client)
)

type nginx struct {
	Title string   `goquery:"h1"`
	Files []string `goquery:"pre a"`
}

type example struct {
	Title string   `goquery:"h1"`
	Files []string `goquery:"table.files tbody tr.js-navigation-item td.content,text"`
}

// GnuUpstreamLatest returns the latest version found from upstream
// mirror. This is still highly experimental
// If the version returned is "0.0.0" then no new versions were found.
func GnuUpstreamLatest(name, url string, current semver.Version) (string, error) {
	var (
		res    *http.Response
		err    error
		latest = "0.0.0"
	)
	if res, err = client.Get(url); err != nil {
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

// ParseName parses a string and returns the name before the last hyphen.
// passing "bash-5.0" would return "bash"
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

// ParseVersion returns the parsed version from a string.
// passing "bash-5.0" should return "5.0"
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

// ParseNameVersion returns the name and version from a string
// the string passed is usually in the GNU form of "name-1.0.0"
func ParseNameVersion(file string) (string, string) {
	return ParseName(file), ParseVersion(file)
}
