package via

import (
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

func TestUpload(t *testing.T) {
	InitClient()
	err := Upload(testDownload)
	checkError(t, err)
}

func TestGetDownloadList(t *testing.T) {
	InitClient()
	_, err := GetDownloadList()
	checkError(t, err)
}

func TestNetRc(t *testing.T) {
	expected := "Mike.Rosset@gmail.com"
	if netrc["login"] != expected {
		t.Errorf("expected %s got %s", expected, netrc["login"])
	}
}

func TestIsUploaded(t *testing.T) {
	pass, err := isUploaded(testDownload)
	if err != nil {
		t.Error(err)
	}
	if !pass {
		t.Errorf("%s is expected to exist on server we got %v",
			testDownload, pass)

	}
	notExpected := "failthis"
	fail, err := isUploaded(notExpected)
	if err != nil {
		t.Error(err)
	}
	if fail {
		t.Errorf("%s is expected to not exist on server we got %v",
			testDownload, fail)
	}
}

func checkError(t *testing.T, err os.Error) {
	if err != nil {
		t.Error(err)
	}
}
