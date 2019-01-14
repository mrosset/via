package plugin

import (
	"github.com/mrosset/via/pkg"
)

type Plugin interface {
	SetConfig(*via.Config)
	Execute() error
}
