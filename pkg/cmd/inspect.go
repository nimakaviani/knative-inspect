package cmd

import (
	"strings"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/flags"
	kiui "github.com/nimakaviani/knative-inspect/pkg/cmd/ui"
	insp "github.com/nimakaviani/knative-inspect/pkg/inspect/ksvc"
	"github.com/spf13/cobra"
)

type InspectCmdOptions struct {
	ui ui.UI

	Debug   bool
	Verbose bool
	opts    insp.InspectorOptions

	kubeconfigFlags flags.KubeconfigFlags
	namespaceFlags  flags.NamespaceFlags
}

func NewInspectCmdOptions(ui *ui.ConfUI, config core.Config, deps core.Deps) *InspectCmdOptions {
	return &InspectCmdOptions{ui: ui, opts: insp.InspectorOptions{Config: config, Deps: deps}}
}

func NewInspectCmd(o *InspectCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Aliases: []string{"i", "insp", "ins"},
		Short:   "Inspect Knative services",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Enable debug output")
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false, "Show verbose output")
	cmd.Flags().StringSliceVarP(&o.opts.Services, "service", "s", []string{}, "Knative services to inspect (can be specified multiple times)")

	o.kubeconfigFlags.Set(cmd)
	o.namespaceFlags.Set(cmd)

	o.opts.Config.ConfigurePathResolver(o.kubeconfigFlags.Path.Value)
	o.opts.Config.ConfigureContextResolver(o.kubeconfigFlags.Context.Value)

	return cmd
}

func (o *InspectCmdOptions) Run() error {
	ui := core.NewPlainUI(o.Debug)
	t1 := time.Now()

	defer func() {
		ui.Debugf("total: %s\n", time.Since(t1))
	}()

	return o.inspect()
}

func (o *InspectCmdOptions) inspect() error {
	result, err := insp.NewInspector(*o.opts.WithNamespace(o.namespaceFlags.Name)).Run()
	if err != nil {
		return err
	}

	kiui.TreeView{Source: strings.Join(o.opts.Services, ","), ResourceMap: result, Verbose: o.Verbose}.Print(o.ui)
	return nil
}
