package tests

import (
	"strings"
	"testing"
)

func TestKubeVersion(t *testing.T) {
	env := BuildEnv(t)
	kni := Kni{t, env.Namespace, Logger{}}

	out, _ := kni.RunWithOpts([]string{"kube-version"}, RunOpts{NoNamespace: true})

	if !strings.Contains(out, "Kubernetes Version: 1.") {
		t.Fatalf("Expected to find kubernetes version")
	}
}
