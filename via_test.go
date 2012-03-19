package via

import (
	"testing"
	"util"
	"util/console"
)

var tests = []string{"ccache"}

var test_plans []*Plan

func init() {
	util.Verbose = false
	for _, t := range tests {
		plan, err := ReadPlan(t)
		util.CheckFatal(err)
		test_plans = append(test_plans, plan)
	}
}

func TestBuildSteps(t *testing.T) {
	for _, plan := range test_plans {
		console.Println("Building", plan)
		if err := BuildSteps(plan); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInstall(t *testing.T) {
	for _, plan := range test_plans {
		console.Println("Installing", plan)
		if err := Install(plan); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRemove(t *testing.T) {
	for _, plan := range test_plans {
		console.Println("Removeing", plan)
		if err := Remove(plan); err != nil {
			t.Fatal(err)
		}
	}
	console.Flush()
}
