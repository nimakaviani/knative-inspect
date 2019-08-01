package cmd

import (
	"fmt"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/fatih/color"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/flags"
	version "github.com/nimakaviani/knative-inspect/pkg/inspect/version"
	"github.com/spf13/cobra"
)

type KubeVersionOptions struct {
	ui    ui.UI
	Debug bool

	kubeconfigFlags flags.KubeconfigFlags

	config core.Config
	deps   core.Deps
}

func NewKubeVersionOptions(ui ui.UI, config core.Config, deps core.Deps) *KubeVersionOptions {
	return &KubeVersionOptions{ui: ui, config: config, deps: deps}
}

func NewKubeVersionCmd(o *KubeVersionOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kube-version",
		Aliases: []string{"kv"},
		Short:   "Print Kubernetes version",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	o.kubeconfigFlags.Set(cmd)
	o.config.ConfigurePathResolver(o.kubeconfigFlags.Path.Value)
	o.config.ConfigureContextResolver(o.kubeconfigFlags.Context.Value)
	return cmd
}

func (o *KubeVersionOptions) Run() error {
	ui := core.NewPlainUI(o.Debug)
	t1 := time.Now()

	defer func() {
		ui.Debugf("total: %s\n", time.Since(t1))
	}()

	return o.version()
}

func (o *KubeVersionOptions) version() error {
	version, err := version.KubeVersion(o.deps)
	if err != nil {
		return err
	}

	o.ui.PrintBlock([]byte(
		fmt.Sprintf(
			"\nKubernetes Version: %s\n",
			color.New(color.FgCyan).Sprintf("%s", version),
		),
	))
	return nil
}
