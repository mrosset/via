package via

import (
	"testing"
	"util"
)

var test = &Plan{}

func init() {
	var err error
	test, err = ReadPlan("ccache")
	util.Verbose = false
	util.CheckFatal(err)
}

func TestDownload(t *testing.T) {
	err := DownloadSrc(test)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStage(t *testing.T) {
	err := Stage(test)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuild(t *testing.T) {
	err := Build(test)
	if err != nil {
		t.Fatal(err)
	}
}

func TestManifest(t *testing.T) {
	_, err := TarManifest(test)
	if err != nil {
		t.Fatal(err)
	}
}

/*
func TestInstall(t *testing.T) {
	err := Install(test.Name)
	if err != nil {
		t.Fatal(err)
	}
}
*/
