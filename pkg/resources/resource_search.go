package resources

import (
	krsc "github.com/nimakaviani/kapp/pkg/kapp/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ResourceSearchOpts struct {
	GroupVersion schema.GroupVersionResource
	ListOpts     metav1.ListOptions

	Namespace string
}

type ResourceSearchImpl struct {
	dynamicClient dynamic.Interface
}

func NewResourceSearch(dynamicClient dynamic.Interface) *ResourceSearchImpl {
	return &ResourceSearchImpl{
		dynamicClient: dynamicClient,
	}
}

func (o *ResourceSearchImpl) Find(opts ResourceSearchOpts) ([]krsc.Resource, error) {
	list, err := o.dynamicClient.Resource(opts.GroupVersion).Namespace(opts.Namespace).List(opts.ListOpts)
	if err != nil {
		return nil, err
	}

	var resources []krsc.Resource
	for _, item := range list.Items {
		resource := krsc.NewResourceUnstructured(item, opts.GroupVersion)
		resources = append(resources, resource)
	}
	return resources, nil
}
