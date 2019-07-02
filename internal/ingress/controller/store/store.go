package store

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Store interface {
	// ListNodes returns a list of all Nodes in the store.
	ListNodes(ctx context.Context) []*corev1.Node
	// ListIngress returns a list of all Ingress objects in the store.
	ListIngress(ctx context.Context) []*extensions.Ingress
}

// Lister contains object listers (stores).
type Listers struct {
	Node    NodeLister
	Ingress IngressLister
}

// Informers defines the required SharedIndexInformers that interact with the API server.
type Informers struct {
}

type kubernetesStore struct {
	client client.Client
	cache  cache.Cache
}

// New creates a new object store to be used in the ingress controller
func New(client client.Client, cache cache.Cache) (Store, error) {
	// ingressInformer, err := cache.GetInformer(&extensions.IngressList{})
	// if err != nil {
	// 	return &kubernetesStore{
	// 		client, cache,
	// 	}, err
	// }

	store := &kubernetesStore{
		client: client,
		cache:  cache,
	}

	//store.listers.Ingress.Store = store.informers.Ingress.GetStore()
	//store.listers.Node.Store = store.informers.Node.GetStore()

	return store, nil
}

func (s kubernetesStore) ListIngress(ctx context.Context) []*extensions.Ingress {
	return nil
}

// ListNodes returns the list of Nodes
func (s kubernetesStore) ListNodes(ctx context.Context) []*corev1.Node {
	var nodes []*corev1.Node
	return nodes
}
