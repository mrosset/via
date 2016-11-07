package via

import (
	"testing"
)

func TestCircuitAddress(t *testing.T) {
	_, err := CircuitAddress()
	if err != nil {
		t.Error(err)
	}
}

func TestCircuit(t *testing.T) {
	err := CircuitBuild("make")
	if err != nil {
		t.Fatal(err)
	}
}
