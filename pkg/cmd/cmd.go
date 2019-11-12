package cmd

import (
	"fmt"
	"io"

	"github.com/cppforlife/cobrautil"
	"github.com/cppforlife/go-cli-ui/ui"
	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	"github.com/spf13/cobra"
)

func NewDefaultCmd(ui *ui.ConfUI) *cobra.Command {
	config := core.NewConfig()
	deps := core.NewDeps(config)

	cmd := &cobra.Command{
		Use:   "kni",
		Short: "Knative Inspect - little debugging tool for Knative services",

		RunE: cobrautil.ShowHelp,

		// Affects children as well
		SilenceErrors: true,
		SilenceUsage:  true,

		// Disable docs header
		DisableAutoGenTag: true,

		// TODO bash completion
	}

	cmd.SetOutput(uiBlockWriter{ui}) // setting output for cmd.Help()

	// TODO bash completion
	cmd.AddCommand(NewInspectCmd(NewInspectCmdOptions(ui, config, deps)))
	cmd.AddCommand(NewConfigChangeCmd(NewConfigChangeOptions(ui, config, deps)))
	cmd.AddCommand(NewKnativeLogsCmd(NewKnativeLogOptions(ui, config, deps)))
	cmd.AddCommand(NewKubeVersionCmd(NewKubeVersionOptions(ui, config, deps)))
	cmd.AddCommand(NewVersionCmd(NewVersionOptions(ui)))

	cobrautil.VisitCommands(cmd, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))

	return cmd
}

func ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return fmt.Errorf("Invalid command - see available commands/subcommands above")
}

type uiBlockWriter struct {
	ui ui.UI
}

var _ io.Writer = uiBlockWriter{}

func (w uiBlockWriter) Write(p []byte) (n int, err error) {
	w.ui.PrintBlock(p)
	return len(p), nil
}
