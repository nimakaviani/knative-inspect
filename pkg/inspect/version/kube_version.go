package inspect

import (
	"fmt"

	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
)

func KubeVersion(deps core.Deps) (string, error) {
	cli, err := deps.DiscoveryClient()
	if err != nil {
		return "", err
	}

	ver, err := cli.ServerVersion()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", ver.Major, ver.Minor), nil
}
