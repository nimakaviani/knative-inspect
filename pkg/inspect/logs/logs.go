package logs

import (
	"fmt"
	"strings"
	"sync"

	klogs "github.com/nimakaviani/kapp/pkg/kapp/logs"
	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LogOptions struct {
	Namespace string
	Follow    bool
	Lines     int64
	Filter    string
	Pods      []string

	Deps   core.Deps
	Config core.Config
}

type Logs struct {
	opts *LogOptions
}

func NewLogs(opts *LogOptions) *Logs {
	return &Logs{opts: opts}
}

func (o *Logs) Run() error {
	err := o.checks()
	if err != nil {
		return err
	}

	coreClient, err := o.opts.Deps.CoreClient()
	if err != nil {
		return err
	}

	pods, err := coreClient.CoreV1().Pods(o.opts.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	lineOpts := (*int64)(nil)
	if o.opts.Lines > 0 {
		lineOpts = &o.opts.Lines
	}

	for _, pod := range pods.Items {
		logOpts := klogs.PodLogOpts{
			Follow:       o.opts.Follow,
			Lines:        lineOpts,
			ContainerTag: true,
		}

		if len(o.opts.Pods) > 0 && !contains(o.opts.Pods, pod.Name) {
			continue
		}

		wg.Add(1)
		pod := pod
		go func() {
			tagFunc := func(cont corev1.Container) string {
				return fmt.Sprintf("%s > %s", pod.Name, cont.Name)
			}

			podClient := coreClient.CoreV1().Pods(o.opts.Namespace)

			dumpUI := NewDumpUI(o.opts.Filter)
			defer dumpUI.Close()

			klogs.NewPodLog(pod, podClient, tagFunc, logOpts).TailAll(dumpUI, make(chan struct{}))
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func contains(slice []string, name string) bool {
	for _, f := range slice {
		if strings.Contains(name, f) {
			return true
		}
	}

	return false
}

func (o Logs) checks() error {
	if o.opts.Filter != "" && !contains([]string{"warn", "info", "error"}, o.opts.Filter) {
		return fmt.Errorf(`--filter "%s" is invalid. It should be one of "error", "info", "warn"`, o.opts.Filter)
	}

	if o.opts.Follow && o.opts.Lines <= 0 {
		return fmt.Errorf("expected --lines to be greater than zero since --follow is not specified")
	}

	return nil
}
