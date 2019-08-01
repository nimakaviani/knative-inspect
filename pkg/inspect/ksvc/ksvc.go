package inspect

import (
	"fmt"
	"log"
	"reflect"

	krsc "github.com/nimakaviani/kapp/pkg/kapp/resources"
	core "github.com/nimakaviani/knative-inspect/pkg/cmd/core"
	resources "github.com/nimakaviani/knative-inspect/pkg/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type KsvcOpts struct {
	Deps core.Deps

	Services  []string
	Namespace string
}

type KSVC struct {
	opts KsvcOpts
}

func NewKsvc(opts KsvcOpts) *KSVC {
	return &KSVC{opts: opts}
}

type ResourceOpts struct {
	kind    string
	group   string
	version string
}

type SearchOpts struct {
	ResourceOpts
	service   string
	name      string
	namespace string

	selectorField string
	selectorValue string
}

func (k *KSVC) Inspect() ([][]krsc.Resource, error) {
	dynamicClient, err := k.opts.Deps.DynamicClient()
	if err != nil {
		return nil, err
	}

	resourceSearch := resources.NewResourceSearch(dynamicClient)
	resourceMap := [][]krsc.Resource{}
	var serviceResource krsc.Resource
	for _, svc := range k.opts.Services {
		ropts := SearchOpts{
			service:   svc,
			name:      svc,
			namespace: k.opts.Namespace,

			ResourceOpts: ResourceOpts{
				kind:    "services",
				group:   "serving.knative.dev",
				version: "v1beta1",
			},

			selectorField: "FieldSelector",
			selectorValue: fmt.Sprintf("metadata.name=%s", svc),
		}

		serviceResource, err = getService(resourceSearch, ropts)
		if err != nil {
			continue
		}

		labeledResources, err := findLabeledResoures(svc, k.opts.Namespace, k.opts.Deps)
		if err != nil {
			return nil, err
		}

		resourceMap = append(resourceMap, append(labeledResources, serviceResource))
	}

	return resourceMap, nil
}

func getService(resourceSearch *resources.ResourceSearchImpl, ropts SearchOpts) (krsc.Resource, error) {
	listOpts := setListOption(&metav1.ListOptions{}, ropts.selectorField, ropts.selectorValue)

	searchOpts := resources.ResourceSearchOpts{
		GroupVersion: schema.GroupVersionResource{
			Group:    ropts.group,
			Version:  ropts.version,
			Resource: ropts.kind,
		},
		ListOpts:  *listOpts,
		Namespace: ropts.namespace,
	}

	rs, err := resourceSearch.Find(searchOpts)
	if err != nil {
		return nil, err
	}

	if len(rs) > 1 {
		panic(fmt.Errorf("found more than one resource for %s:%s/%s:%s", ropts.group, ropts.version, ropts.kind, ropts.name))
	}

	if len(rs) == 0 {
		log.Fatalf("service %s not found", ropts.service)
	}

	return rs[0], nil
}

func setListOption(opts *metav1.ListOptions, fieldName string, fieldValue string) *metav1.ListOptions {
	ps := reflect.ValueOf(opts)
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		f := s.FieldByName(fieldName)
		if f.IsValid() {
			if f.CanSet() {
				if f.Kind() == reflect.String {
					f.SetString(fieldValue)
				}
			}
		}
	}

	return opts
}

func findLabeledResoures(service, namespace string, deps core.Deps) ([]krsc.Resource, error) {

	filterLabel := map[string]string{
		"serving.knative.dev/service": service,
	}

	dynamicClient, err := deps.DynamicClient()
	if err != nil {
		return nil, err
	}

	coreClient, err := deps.CoreClient()
	if err != nil {
		return nil, err
	}

	idr := krsc.NewIdentifiedResources(coreClient, dynamicClient, []string{})
	return idr.List(labels.Set(filterLabel).AsSelector())
}
