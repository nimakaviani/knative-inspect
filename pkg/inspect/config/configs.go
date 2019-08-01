package config

import (
	"fmt"
	"log"
	"reflect"
	"time"

	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// TODO: look for reason - prevent config reflection from panicking
	_ "knative.dev/pkg/system/testing"

	diff "github.com/walmartlabs/object-diff/pkg/obj_diff"

	// "knative.dev/pkg/configmap"
	apisconfig "knative.dev/serving/pkg/apis/config"
	autoscalerconfig "knative.dev/serving/pkg/autoscaler"
	deployment "knative.dev/serving/pkg/deployment"
	gcconfig "knative.dev/serving/pkg/gc"
	networkconfig "knative.dev/serving/pkg/network"
	certconfig "knative.dev/serving/pkg/reconciler/certificate/config"
	istioconfig "knative.dev/serving/pkg/reconciler/ingress/config"
	domainconfig "knative.dev/serving/pkg/reconciler/route/config"
	tracingconfig "knative.dev/serving/pkg/tracing/config"
)

type ConfigMapChangeOpts struct {
	Namespace      string
	ConfigMapNames []string

	Deps   core.Deps
	Config core.Config
}

type ConfigInspector struct {
	opts           ConfigMapChangeOpts
	skippedNames   map[string]interface{}
	knativeConfigs map[string]interface{}
}

func NewConfigInspector(opts ConfigMapChangeOpts) *ConfigInspector {
	logger := &logger{}
	d, _ := time.ParseDuration("1m")

	knativeConfigs := map[string]interface{}{
		tracingconfig.ConfigName:         tracingconfig.NewTracingConfigFromConfigMap,
		apisconfig.DefaultsConfigName:    apisconfig.NewDefaultsConfigFromConfigMap,
		gcconfig.ConfigName:              gcconfig.NewConfigFromConfigMapFunc(logger, d),
		domainconfig.DomainConfigName:    domainconfig.NewDomainFromConfigMap,
		networkconfig.ConfigName:         networkconfig.NewConfigFromConfigMap,
		certconfig.CertManagerConfigName: certconfig.NewCertManagerConfigFromConfigMap,
		istioconfig.IstioConfigName:      istioconfig.NewIstioFromConfigMap,
		autoscalerconfig.ConfigName:      autoscalerconfig.NewConfigFromConfigMap,
		deployment.ConfigName:            deployment.NewConfigFromConfigMap,
	}

	skippedNames := make(map[string]interface{})
	if len(opts.ConfigMapNames) > 0 {
		for k, _ := range knativeConfigs {
			found := false
			for _, n := range opts.ConfigMapNames {
				if k == n {
					found = true
					break
				}
			}
			if !found {
				skippedNames[k] = struct{}{}
			}
		}
	}

	return &ConfigInspector{
		opts:           opts,
		skippedNames:   skippedNames,
		knativeConfigs: knativeConfigs,
	}
}

func (o *ConfigInspector) Run() (map[string]*diff.ChangeSet, error) {
	coreClient, err := o.opts.Deps.CoreClient()
	if err != nil {
		return nil, err
	}

	configs, err := coreClient.CoreV1().ConfigMaps(o.opts.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	changes := make(map[string]*diff.ChangeSet)

	for _, config := range configs.Items {
		if fn, ok := o.knativeConfigs[config.Name]; ok {

			if _, ok := o.skippedNames[config.Name]; ok {
				continue
			}

			defaultConfig, err := call(fn, &corev1.ConfigMap{})
			if err != nil {
				log.Printf("default config for %s => %s", config.Name, err.Error())
				continue
			}

			appliedConfig, err := call(fn, &config)
			if err != nil {
				log.Printf("applied config for %s => %s", config.Name, err.Error())
				continue
			}

			diffs, err := diff.Diff(defaultConfig, appliedConfig)
			if err != nil {
				log.Printf("diff failed for %s", config.Name)
				continue
			}

			changes[config.Name] = diffs
		}
	}

	return changes, nil
}

func call(fn interface{}, configMap *corev1.ConfigMap) (interface{}, error) {
	constructor := reflect.ValueOf(fn)
	inputs := []reflect.Value{reflect.ValueOf(configMap)}
	outputs := constructor.Call(inputs)

	result := outputs[0].Interface()
	errVal := outputs[1]

	if !errVal.IsNil() {
		return nil, fmt.Errorf("%q", errVal.Interface())
	}

	return result, nil
}
