package store

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

// NodeLister makes a Store that lists Nodes.
type NodeLister struct {
	cache.Store
}

// ByKey returns the Node matching key in the local Node Store.
func (nl *NodeLister) ByKey(key string) (*apiv1.Node, error) {
	n, exists, err := nl.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, err
	}
	return n.(*apiv1.Node), nil
}
