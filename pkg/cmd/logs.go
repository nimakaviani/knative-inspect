package cmd

import (
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	"github.com/nimakaviani/knative-inspect/pkg/cmd/flags"
	"github.com/nimakaviani/knative-inspect/pkg/inspect/logs"
	"github.com/spf13/cobra"
)

type KnativeLogOptions struct {
	ui ui.UI

	Debug bool
	opts  *logs.LogOptions

	kubeconfigFlags flags.KubeconfigFlags
	namespaceFlags  flags.NamespaceFlags
}

func NewKnativeLogOptions(ui *ui.ConfUI, config core.Config, deps core.Deps) *KnativeLogOptions {
	return &KnativeLogOptions{
		ui:   ui,
		opts: &logs.LogOptions{Config: config, Deps: deps},
	}
}

func NewKnativeLogsCmd(o *KnativeLogOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"l"},
		Short:   "Knative Component Logs",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	cmd.Flags().BoolVarP(&o.opts.Follow, "follow", "f", false, "As new pods are added, new pod logs will be printed")
	cmd.Flags().Int64Var(&o.opts.Lines, "lines", 10, "Limit to number of lines (use -1 to remove limit)")
	cmd.Flags().StringVar(&o.opts.Filter, "filter", "", "Filter logs by log type (info, warn, error)")
	cmd.Flags().StringSliceVarP(
		&o.opts.Pods,
		"pods", "p",
		[]string{},
		"Knative component to grab logs for (can be specified multiple times, eg. activator, autoscaler)",
	)

	o.kubeconfigFlags.Set(cmd)
	o.namespaceFlags.Set(cmd, servingNamespace)

	o.opts.Config.ConfigurePathResolver(o.kubeconfigFlags.Path.Value)
	o.opts.Config.ConfigureContextResolver(o.kubeconfigFlags.Context.Value)
	return cmd
}

func (o *KnativeLogOptions) Run() error {
	ui := core.NewPlainUI(o.Debug)
	t1 := time.Now()

	defer func() {
		ui.Debugf("total: %s\n", time.Since(t1))
	}()

	return o.logs()
}

func (o *KnativeLogOptions) logs() error {
	// set the name based on whatever comes through from the flag
	o.opts.Namespace = o.namespaceFlags.Name

	return logs.NewLogs(o.opts).Run()
}
