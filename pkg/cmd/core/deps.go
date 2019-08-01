package core

import (
	"fmt"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Deps interface {
	DynamicClient() (dynamic.Interface, error)
	CoreClient() (kubernetes.Interface, error)
	DiscoveryClient() (discovery.DiscoveryInterface, error)
}

type DepsImpl struct {
	config Config
}

var _ Deps = &DepsImpl{}

func NewDeps(config Config) Deps {
	return &DepsImpl{config}
}

func (f *DepsImpl) DynamicClient() (dynamic.Interface, error) {
	config, err := f.config.RESTConfig()
	if err != nil {
		return nil, err
	}

	// TODO high QPS
	config.QPS = 1000
	config.Burst = 1000

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("building Dynamic clientset: %s", err)
	}

	return clientset, nil
}

func (f *DepsImpl) CoreClient() (kubernetes.Interface, error) {
	config, err := f.config.RESTConfig()
	if err != nil {
		return nil, err
	}

	config.QPS = 1000
	config.Burst = 1000

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("building Core clientset: %s", err)
	}

	return clientset, nil
}

func (f *DepsImpl) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	config, err := f.config.RESTConfig()
	if err != nil {
		return nil, err
	}

	config.QPS = 1000
	config.Burst = 1000

	clientset, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("building discovery clientset: %s", err)
	}

	return clientset, nil
}
