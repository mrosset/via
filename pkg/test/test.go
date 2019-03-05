package test

import (
	"reflect"
	"testing"
)

//ExpectGotFmt format for failed test output
const ExpectGotFmt = "%s: expect '%v' got '%v'"

// Test is a type for a test
type Test struct {
	Name     string
	Expect   interface{}
	Got      interface{}
	GotFn    func() (interface{}, error)
	GotArgFn func(interface{}) (interface{}, error)
}

// Tests is a slice of tests
type Tests []Test

// Equals runs assertions on each Test
func (ts Tests) Equals(t *testing.T) {
	for _, test := range ts {
		test.Equals(t)
	}
}

// Equals asserts the Test
func (vt Test) Equals(t *testing.T) bool {
	var (
		err error
	)
	if vt.GotFn != nil {
		vt.Got, err = vt.GotFn()
		if err != nil {
			t.Errorf(ExpectGotFmt, vt.Name, nil, err)
			return false
		}
	}
	if !reflect.DeepEqual(vt.Expect, vt.Got) {
		t.Errorf(ExpectGotFmt, vt.Name, vt.Expect, vt.Got)
		return false
	}
	return true
}
