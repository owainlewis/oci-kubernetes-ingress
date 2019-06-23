package store

import (
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Store interface {
	// ListNodes returns a list of all Nodes in the store.
	ListNodes() []*corev1.Node
}

// Lister contains object listers (stores).
type Lister struct {
	Node NodeLister
}

// Informer defines the required SharedIndexInformers that interact with the API server.
type Informer struct {
	Node cache.SharedIndexInformer
}

type kubernetesStore struct {
	listers   *Lister
	informers *Informer
}

// New creates a new object store to be used in the ingress controller
func New(mgr manager.Manager) (Store, error) {
	store := &kubernetesStore{
		listers: &Lister{},
	}

	mgrCache := mgr.GetCache()
	var err error
	store.informers.Node, err = mgrCache.GetInformer(&corev1.Node{})
	if err != nil {
		return nil, err
	}
	store.listers.Node.Store = store.informers.Node.GetStore()

	return store, nil
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
