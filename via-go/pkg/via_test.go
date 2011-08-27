package via

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"github.com/kr/pretty.go"
)

func testListPlans(t *testing.T) {
	err := ListPlans()
	checkError(t, err)
}

func TestFindPlan(t *testing.T) {
	expected := "bash"
	plan, err := FindPlan(expected)
	checkError(t, err)
	if plan.Name != expected {
		t.Errorf("exected %s for Name got %s", expected, plan.Name)
	}
}

func TestPackage(t *testing.T) {
	err := Package("bash", "x86_64")
	checkError(t, err)
}

func testHeaders(t *testing.T) {
	plan, err := FindPlan("vim")
	checkError(t, err)
	file := filepath.Join(repo, "x86_64", plan.NameVersion()+"-x86_64.tar.gz")
	//file := "/home/strings/via/cache/packages/vim-7.3.tar.gz"
	tbr, err := NewTarBallReader(file)
	checkError(t, err)
	for {
		hdr, err := tbr.tr.Next()
		if err == os.EOF {
			break
		}
		if hdr.Typeflag == tar.TypeSymlink {
			fmt.Printf("%# v", pretty.Formatter(hdr))
		}
	}
}

func TestUnPack(t *testing.T) {
	plan, err := FindPlan("bash")
	checkError(t, err)
	file := filepath.Join(repo, "x86_64", plan.NameVersion()+"-x86_64.tar.gz")
	err = Unpack("./tmp", file)
	checkError(t, err)
}

func checkError(t *testing.T, err os.Error) {
	if err != nil {
		t.Error(err)
	}
}
