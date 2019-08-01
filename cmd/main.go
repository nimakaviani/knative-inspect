package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/nimakaviani/knative-inspect/pkg/cmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()

	command := cmd.NewDefaultCmd(confUI)

	err := command.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
