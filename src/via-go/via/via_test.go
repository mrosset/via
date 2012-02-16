package via

import (
	"fmt"
	"os"
	"testing"
)

var (
	//tests    = []string{"bash", "ncdu", "file", "coreutils", "eglibc", "git"}
	tests    = []string{"ncdu", "bash"}
	testArch = "x86_64"
	testRoot = "./tmp"
)

func init() {
	os.Mkdir(testRoot, 0755)
}

func TestFindPlan(t *testing.T) {
	for _, test := range tests {
		expected := test
		plan, err := FindPlan(expected)
		if err != nil {
			t.Error(err)
		}
		if plan.Name != expected {
			t.Errorf("exected %s for Name got %s", expected, plan.Name)
		}
	}
}

func TestPackage(t *testing.T) {
	for _, test := range tests {
		err := Package(test, testArch)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInstall(t *testing.T) {
	for _, test := range tests {
		err := Install(testRoot, test)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateRepo(t *testing.T) {
	err := UpdateRepo(testArch)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadRepo(t *testing.T) {
	_, err := LoadRepo(testArch)
	if err != nil {
		t.Error(err)
	}
}

func TestNetRc(t *testing.T) {
	expected := "Mike.Rosset@gmail.com"
	if netrc["login"] != expected {
		t.Errorf("expected %s got %s", expected, netrc["login"])
	}
}

func testUploadRepo(t *testing.T) {
	if err := uploadRepo(testArch); err != nil {
		t.Error(err)
	}
}

func testGetDownloadList(t *testing.T) {
	list, err := GetDownloadList()
	if err != nil {
		t.Error(err)
	}
	for i, l := range list {
		fmt.Printf("%-0.2d %s\n", i, l)
	}
}

func TestDownloadSrc(t *testing.T) {
	url := "http://mirrors.kernel.org/gnu/bash/bash-4.2.tar.gz"
	if err := DownloadSrc(url); err != nil {
		t.Error(err)
	}
}

func TestDownloadSig(t *testing.T) {
	url := "http://mirrors.kernel.org/gnu/bash/bash-4.2.tar.gz"
	if err := DownloadSig(url); err != nil {
		t.Error(err)
	}

}

func TestCheck(t *testing.T) {
	for _, test := range tests {
		err := Check(testRoot, test)
		if err != nil {
			t.Error(err)
		}
	}
}

func testOwnsFile(t *testing.T) {
	expected := "file"
	mani, err := OwnsFile(testRoot, "libmagic.so.1")
	if err != nil {
		t.Error(err)
	}
	if mani.Meta.Name != expected {
		t.Errorf("expected %s got %s", expected, mani.Meta.Name)
	}
}

func testRemove(t *testing.T) {
	for _, test := range tests {
		err := Remove(testRoot, test)
		if err != nil {
			t.Error(err)
		}
	}
}
