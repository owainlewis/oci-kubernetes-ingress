package ingress

import (
	"errors"

	"github.com/golang/glog"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/config"
	"k8s.io/api/extensions/v1beta1"
)

// Manager transforms ingress objects into OCI load balancer specifications
// and ensures that the desired and actual state of the world align.
type Manager interface {
	EnsureIngress(specification Specification) error
	EnsureIngressDeleted(ingress *v1beta1.Ingress) error
}

type defaultManager struct {
	conf config.Config
	svc  *LoadBalancerService
}

// NewDefaultManager constructs a new ingress Manager
func NewDefaultManager(conf config.Config) (Manager, error) {
	svc, err := NewLoadBalancerService(conf)
	if err != nil {
		return nil, err
	}
	return &defaultManager{
		conf: conf,
		svc:  svc,
	}, nil
}

// EnsureIngress will ensure that the observed Ingress object state is reflected
// in OCI
func (mgr *defaultManager) EnsureIngress(specification Specification) error {
	glog.V(4).Infof("Ensuring ingress for specification: %s", specification.Ingress.Name)
	// Check if a load balancer exists already for this ingress object
	_, err := mgr.svc.GetLoadBalancer(specification.Config.Loadbalancer.Compartment, GetLoadBalancerUniqueName(specification.Ingress))
	// We cannot find a load balancer for this ingress or something went wrong.
	if err != nil {
		// If no load balancer exists then create one
		glog.V(4).Infof("Creating a new loadbalancer for specification: %s", specification.Ingress.Name)
		_, err = mgr.svc.CreateLoadBalancer(specification)
		return err
	}

	// If the load balancer exists update it
	glog.V(4).Infof("Updating load balancer for specification: %s", specification.Ingress.Name)

	return nil
}

// EnsureIngressDeleted will ensure that an ingress object is removed from OCI
func (mgr *defaultManager) EnsureIngressDeleted(ingress *v1beta1.Ingress) error {
	if ingress == nil {
		return errors.New("Trying to delete ingress which is nil")
	}

	name := GetLoadBalancerUniqueName(ingress)

	glog.Infof("Deleting ingress %s", name)
	return mgr.svc.DeleteLoadBalancer(mgr.conf.Loadbalancer.Compartment, name)
}
