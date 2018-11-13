package ingress

import (
	"errors"

	"github.com/golang/glog"
	"github.com/owainlewis/oci-kubernetes-ingress/pkg/config"
	"k8s.io/api/extensions/v1beta1"
)

// Manager transforms ingress objects into OCI load balancer specifications
// and ensures that the desired and actual state of the world align.
type Manager interface {
	EnsureIngress(specification Specification) error
	EnsureIngressDeleted(ingress *v1beta1.Ingress) error
}

type defaultManager struct {
	svc *LoadBalancerService
}

// NewDefaultManager constructs a new ingress Manager
func NewDefaultManager(conf config.Config) (Manager, error) {
	svc, err := NewLoadBalancerService(conf)
	if err != nil {
		return nil, err
	}
	return &defaultManager{
		svc: svc,
	}, nil
}

// EnsureIngress will ensure that the observed Ingress object state is reflected
// in OCI
func (mgr *defaultManager) EnsureIngress(specification Specification) error {
	glog.Infof("Ensuring ingress")
	glog.Infof("Ingress Spec: %+v", specification)

	// Check if a load balancer exists already for this ingress object
	_, err := mgr.svc.GetLoadBalancer(specification)
	if err != nil {
		return err
	}

	// If the load balancer exists update it if needed

	// If no load balancer exists then create one
	_, err = mgr.svc.CreateLoadBalancer(specification)
	return err
}

// EnsureIngressDeleted will ensure that an ingress object is removed from OCI
func (mgr *defaultManager) EnsureIngressDeleted(ingress *v1beta1.Ingress) error {
	if ingress == nil {
		return errors.New("Trying to delete ingress which is nil")
	}

	name := GetLoadBalancerUniqueName(ingress)

	glog.Infof("Deleting ingress %s", name)
	return mgr.svc.DeleteLoadBalancer(name)
}
