package tests

import (
	// "reflect"
	"strings"
	"testing"
	"time"
)

func TestInspect(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kni := Kni{t, env.Namespace, logger}
	kubectl := Kubectl{t, env.Namespace, logger}

	yaml := `
apiVersion: serving.knative.dev/v1alpha1
kind: Service
metadata:
  name: testapp
spec:
  template:
    spec:
      containers:
      - image: nimak/knative-sample-app:v3
        env:
        - name: NAME
          value: "John"
`

	cleanUp := func() {
		kubectl.RunWithOpts([]string{"delete", "-f", "-"}, RunOpts{AllowError: true, StdinReader: strings.NewReader(yaml)})
	}

	cleanUp()
	defer cleanUp()

	logger.Section("deploy initial", func() {
		kubectl.RunWithOpts([]string{"apply", "-f", "-"}, RunOpts{StdinReader: strings.NewReader(yaml)})
	})

	logger.Section("check resources", func() {

		time.Sleep(1 * time.Second)

		out, _ := kni.RunWithOpts([]string{"inspect", "-s", "testapp"}, RunOpts{})
		bufferReader := BufferReader(strings.NewReader(out))

		resources := []string{
			"Service",
			"Configuration",
			"Revision",
			"Deployment",
			"ReplicaSet",
			"Image",
			"PodAutoscaler",
			"Service",
			"ServerlessService",
		}

		for _, r := range resources {
			bufferReader.ShouldSay("testapp")
			bufferReader.ShouldSay(r)
		}
	})
}
