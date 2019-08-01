package flags

import (
	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
)

type FlagsFactory struct {
	config core.Config
	deps   core.Deps
}

func NewFlagsFactory(config core.Config, depconfig core.Deps) FlagsFactory {
	return FlagsFactory{config, depconfig}
}
