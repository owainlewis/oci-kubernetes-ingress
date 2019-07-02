package store

import (
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

// IngressLister makes a Store that lists Ingress.
type IngressLister struct {
	cache.Store
}

// ByKey returns an Ingress by key
func (il IngressLister) ByKey(key string) (*extensions.Ingress, error) {
	i, exists, err := il.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, err
	}
	return i.(*extensions.Ingress), nil
}
