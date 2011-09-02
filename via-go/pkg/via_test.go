package via

import (
	"fmt"
	"os"
	"testing"
)

var (
	tests    = []string{"bash", "ncdu"}
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
		checkError(t, err)
		if plan.Name != expected {
			t.Errorf("exected %s for Name got %s", expected, plan.Name)
		}
	}
}

func TestPackage(t *testing.T) {
	for _, test := range tests {
		err := Package(test, testArch)
		checkError(t, err)
	}
}

func TestUnPack(t *testing.T) {
	for _, test := range tests {
		plan, err := FindPlan(test)
		checkError(t, err)
		err = Unpack(testRoot, PkgAbsFile(plan, testArch))
		checkError(t, err)
	}
}

func TestUpdateRepo(t *testing.T) {
	err := UpdateRepo(testArch)
	checkError(t, err)
}

func TestLoadRepo(t *testing.T) {
	_, err := LoadRepo(testArch)
	checkError(t, err)
}

var testDownload = "bash-4.2-x86_64.tar.bz2"

func TestDownload(t *testing.T) {
	InitClient()
	err := Download(testDownload)
	checkError(t, err)
}

func testUpload(t *testing.T) {
	InitClient()
	err := upload(testDownload)
	checkError(t, err)
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

func TestGetDownloadList(t *testing.T) {
	InitClient()
	list, err := GetDownloadList()
	for i, l := range list {
		fmt.Printf("%-0.2d %s\n", i, l)
	}
	checkError(t, err)
}

func checkError(t *testing.T, err os.Error) {
	if err != nil {
		t.Error(err)
	}
}
