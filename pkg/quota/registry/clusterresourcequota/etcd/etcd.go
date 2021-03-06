package etcd

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	kapi "k8s.io/kubernetes/pkg/api"

	"github.com/openshift/origin/pkg/quota/api"
	"github.com/openshift/origin/pkg/quota/registry/clusterresourcequota"
	"github.com/openshift/origin/pkg/util/restoptions"
)

type REST struct {
	*registry.Store
}

// NewREST returns a RESTStorage object that will work against ClusterResourceQuota objects.
func NewREST(optsGetter restoptions.Getter) (*REST, *StatusREST, error) {
	store := &registry.Store{
		Copier:            kapi.Scheme,
		NewFunc:           func() runtime.Object { return &api.ClusterResourceQuota{} },
		NewListFunc:       func() runtime.Object { return &api.ClusterResourceQuotaList{} },
		PredicateFunc:     clusterresourcequota.Matcher,
		QualifiedResource: api.Resource("clusterresourcequotas"),

		CreateStrategy: clusterresourcequota.Strategy,
		UpdateStrategy: clusterresourcequota.Strategy,
		DeleteStrategy: clusterresourcequota.Strategy,
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: clusterresourcequota.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, nil, err
	}

	statusStore := *store
	statusStore.CreateStrategy = nil
	statusStore.DeleteStrategy = nil
	statusStore.UpdateStrategy = clusterresourcequota.StatusStrategy

	return &REST{store}, &StatusREST{store: &statusStore}, nil
}

// StatusREST implements the REST endpoint for changing the status of a resourcequota.
type StatusREST struct {
	store *registry.Store
}

// StatusREST implements Patcher
var _ = rest.Patcher(&StatusREST{})

func (r *StatusREST) New() runtime.Object {
	return &api.ClusterResourceQuota{}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx apirequest.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx apirequest.Context, name string, objInfo rest.UpdatedObjectInfo) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo)
}
