package plugin

import (
	"github.com/mrosset/via/pkg"
	"testing"
)

func testPluginBuild(t *testing.T) {
	config := via.GetConfig()
	err := Build(config)
	if err != nil {
		t.Error(err)
	}
}
