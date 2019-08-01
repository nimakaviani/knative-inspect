package cmd

import (
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/flags"
	kiui "github.com/nimakaviani/knative-inspect/pkg/cmd/ui"
	cfg "github.com/nimakaviani/knative-inspect/pkg/inspect/config"
	"github.com/spf13/cobra"
)

const (
	servingNamespace = "knative-serving"
)

type ConfigChangeOptions struct {
	ui ui.UI

	Debug bool
	opts  cfg.ConfigMapChangeOpts

	kubeconfigFlags flags.KubeconfigFlags
	namespaceFlags  flags.NamespaceFlags
}

func NewConfigChangeOptions(ui *ui.ConfUI, config core.Config, deps core.Deps) *ConfigChangeOptions {
	return &ConfigChangeOptions{
		ui:   ui,
		opts: cfg.ConfigMapChangeOpts{Config: config, Deps: deps},
	}
}

func NewConfigChangeCmd(o *ConfigChangeOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config-changes",
		Aliases: []string{"cc"},
		Short:   "Changes to Knative ConfigMaps",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	cmd.Flags().StringSliceVarP(
		&o.opts.ConfigMapNames,
		"config", "c",
		[]string{},
		"Knative config maps to check (can be specified multiple times)",
	)

	o.kubeconfigFlags.Set(cmd)
	o.namespaceFlags.Set(cmd, servingNamespace)

	o.opts.Config.ConfigurePathResolver(o.kubeconfigFlags.Path.Value)
	o.opts.Config.ConfigureContextResolver(o.kubeconfigFlags.Context.Value)
	return cmd
}

func (o *ConfigChangeOptions) Run() error {
	ui := core.NewPlainUI(o.Debug)
	t1 := time.Now()

	defer func() {
		ui.Debugf("total: %s\n", time.Since(t1))
	}()

	return o.changes()
}

func (o *ConfigChangeOptions) changes() error {
	// set the name based on whatever comes through from the flag
	o.opts.Namespace = o.namespaceFlags.Name

	diffs, err := cfg.NewConfigInspector(o.opts).Run()
	if err != nil {
		return err
	}

	kiui.DiffView{Source: "Knative ConfigMaps", ChangeSet: diffs}.Print(o.ui)
	return nil
}
