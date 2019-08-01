package cmd

import (
	"fmt"
	"io"
	"strings"

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

	// Last one runs first
	cobrautil.VisitCommands(cmd, reconfigureCmdWithSubcmd)
	cobrautil.VisitCommands(cmd, reconfigureLeafCmd)
	cobrautil.VisitCommands(cmd, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))

	return cmd
}

func reconfigureCmdWithSubcmd(cmd *cobra.Command) {
	if len(cmd.Commands()) == 0 {
		return
	}

	if cmd.Args == nil {
		cmd.Args = cobra.ArbitraryArgs
	}
	if cmd.RunE == nil {
		cmd.RunE = ShowSubcommands
	}

	var strs []string
	for _, subcmd := range cmd.Commands() {
		strs = append(strs, subcmd.Use)
	}

	cmd.Short += " (" + strings.Join(strs, ", ") + ")"
}

func reconfigureLeafCmd(cmd *cobra.Command) {
	if len(cmd.Commands()) > 0 {
		return
	}

	if cmd.RunE == nil {
		panic(fmt.Sprintf("Internal: Command '%s' does not set RunE", cmd.CommandPath()))
	}

	if cmd.Args == nil {
		origRunE := cmd.RunE
		cmd.RunE = func(cmd2 *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("command '%s' does not accept extra arguments '%s'", args[0], cmd2.CommandPath())
			}
			return origRunE(cmd2, args)
		}
		cmd.Args = cobra.ArbitraryArgs
	}
}

func ShowSubcommands(cmd *cobra.Command, args []string) error {
	var strs []string
	for _, subcmd := range cmd.Commands() {
		strs = append(strs, subcmd.Use)
	}
	return fmt.Errorf("Use one of available subcommands: %s", strings.Join(strs, ", "))
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
