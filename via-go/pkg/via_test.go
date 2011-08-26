package via

import (
	"os"
	"testing"
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

func checkError(t *testing.T, err os.Error) {
	if err != nil {
		t.Error(err)
	}
}
