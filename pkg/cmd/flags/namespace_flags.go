package flags

import (
	"github.com/spf13/cobra"
)

type NamespaceFlags struct {
	Name string
}

func (s *NamespaceFlags) Set(cmd *cobra.Command, ns ...string) {
	if len(ns) > 1 {
		panic("cannot have multiple namespaces")
	}

	defaultNS := "default"
	if len(ns) == 1 {
		defaultNS = ns[0]
	}

	cmd.Flags().StringVarP(&s.Name, "namespace", "n", defaultNS, "Knative service namespace")
}
