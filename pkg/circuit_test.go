package via

import (
	"testing"
)

func testCircuitAddress(t *testing.T) {
	_, err := CircuitAddress()
	if err != nil {
		t.Error(err)
	}
}

func testCircuit(t *testing.T) {
	err := CircuitBuild("make")
	if err != nil {
		t.Fatal(err)
	}
}
