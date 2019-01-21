package via

import (
	"testing"
)

func TestContainer(t *testing.T) {
	Container()
	t.Errorf("this should fail")
}
