package via

import (
	"testing"
	"util"
)

var bash = &Plan{}

func init() {
	util.CheckFatal(ReadJson(bash, config.Plans()+"/bash.json"))
}

func TestDownload(t *testing.T) {
	err := DownloadSrc(bash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStage(t *testing.T) {
	err := Stage(bash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuild(t *testing.T) {
	err := Build(bash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstall(t *testing.T) {
	err := Install(bash)
	if err != nil {
		t.Fatal(err)
	}
}
