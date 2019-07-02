package store

import (
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"

	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Store interface {
	// ListNodes returns a list of all Nodes in the store.
	ListNodes() []*corev1.Node
	ListIngress() []*extensions.Ingress
}

// Lister contains object listers (stores).
type Listers struct {
	Node    NodeLister
	Ingress IngressLister
}

// Informer defines the required SharedIndexInformers that interact with the API server.
type Informers struct {
	Ingress cache.SharedIndexInformer
	Node    cache.SharedIndexInformer
}

type kubernetesStore struct {
	listers   *Listers
	informers *Informers
}

// New creates a new object store to be used in the ingress controller
func New(mgr manager.Manager) (Store, error) {
	store := &kubernetesStore{
		informers: &Informers{},
		listers:   &Listers{},
	}

	store.listers.Ingress.Store = store.informers.Ingress.GetStore()
	store.listers.Node.Store = store.informers.Node.GetStore()

	return store, nil
}

func (s kubernetesStore) ListIngress() []*extensions.Ingress {
	var ingress []*extensions.Ingress
	for _, item := range s.listers.Ingress.List() {
		n := item.(*extensions.Ingress)
		ingress = append(ingress, n)
	}

	return ingress
}

// ListNodes returns the list of Nodes
func (s kubernetesStore) ListNodes() []*corev1.Node {
	var nodes []*corev1.Node
	for _, item := range s.listers.Node.List() {
		n := item.(*corev1.Node)
		nodes = append(nodes, n)
	}

	return nodes
}
