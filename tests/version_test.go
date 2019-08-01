package tests

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	env := BuildEnv(t)
	kni := Kni{t, env.Namespace, Logger{}}

	out, _ := kni.RunWithOpts([]string{"version"}, RunOpts{NoNamespace: true})

	if !strings.Contains(out, "Client Version") {
		t.Fatalf("Expected to find client version")
	}
}
