package ui

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/fatih/color"
	diff "github.com/walmartlabs/object-diff/pkg/obj_diff"
)

type DiffView struct {
	Source    string
	ChangeSet map[string]*diff.ChangeSet
}

func (v DiffView) Print(ui ui.UI) {
	var result []byte

	ui.PrintLinef("Changes in %s", v.Source)

	for name, changeSet := range v.ChangeSet {
		if len(changeSet.Changes) <= 0 {
			continue
		}

		result = append(
			result,
			[]byte(color.New(color.FgYellow).Sprintf("* %s\n", name))...,
		)

		for _, c := range changeSet.Changes {
			if c.IsDeletion() {
				result = append(
					result,
					[]byte(color.New(color.FgRed).Sprintf("- %v\n", c.PathString()))...,
				)
			}

			if c.IsAddition() {
				result = append(
					result,
					[]byte(color.New(color.FgGreen).Sprintf("+ %v %v\n", c.PathString(), c.GetNewValue()))...,
				)
			}

			if !c.IsDeletion() && !c.IsAddition() && c.GetNewValue() != c.GetOldValue() {
				result = append(
					result,
					[]byte(color.New(color.FgRed).Sprintf("- %v %v\n", c.PathString(), c.GetOldValue()))...,
				)
				result = append(
					result,
					[]byte(color.New(color.FgGreen).Sprintf("+ %v %v\n", c.PathString(), c.GetNewValue()))...,
				)
			}
		}
		result = append(result, []byte("\n")...)
	}

	ui.PrintBlock(result)
}
