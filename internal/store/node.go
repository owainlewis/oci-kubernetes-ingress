package store

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// A NodeStore is a thin abstraction over a cache.Store
// with logic built in to deal with type casting etc.
type NodeStore struct {
	cache.Store
}

// NewNodeStoreFromCache will return a custom store that knows about K8s node objects.
func NewNodeStoreFromCache(cache cache.Store) NodeStore {
	return NodeStore{cache}
}

// ListNodes will return a list of nodes from the cache store
func (s NodeStore) ListNodes() []*corev1.Node {
	nodes := []*corev1.Node{}
	for _, item := range s.List() {
		node := item.(*corev1.Node)
		nodes = append(nodes, node)
	}

	return nodes
}
