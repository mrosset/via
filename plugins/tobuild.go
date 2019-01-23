package main

import (
	"github.com/mrosset/util/console"
	"github.com/mrosset/via/pkg"
)

type tobuild struct {
	config *via.Config
}

func (t *tobuild) SetConfig(config *via.Config) {
	t.config = config
}

func (t *tobuild) Execute() error {
	plans, err := via.GetPlans()
	if err != nil {
		return err
	}
	for _, p := range plans {
		if !p.IsRebuilt {
			console.Println(p.NameVersion(), p.IsRebuilt)
		}
	}
	console.Flush()
	return nil
}

var Tobuild tobuild
