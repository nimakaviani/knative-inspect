package inspect

import (
	krsc "github.com/nimakaviani/kapp/pkg/kapp/resources"
	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
)

type ResourceMap [][]krsc.Resource

type Inspector interface {
	Run() (ResourceMap, error)
}

type InspectorOptions struct {
	Namespace string
	Services  []string

	Config core.Config
	Deps   core.Deps
}

func (k *InspectorOptions) WithNamespace(namespace string) *InspectorOptions {
	k.Namespace = namespace
	return k
}

type inspectorImpl struct {
	opts InspectorOptions
}

func NewInspector(opts InspectorOptions) Inspector {
	return &inspectorImpl{opts}
}

func (o *inspectorImpl) Run() (ResourceMap, error) {
	ksvcOpts := KsvcOpts{
		Services:  o.opts.Services,
		Namespace: o.opts.Namespace,
		Deps:      o.opts.Deps,
	}

	resourceMap, err := NewKsvc(ksvcOpts).Inspect()
	if err != nil {
		return nil, err
	}

	return resourceMap, nil
}
