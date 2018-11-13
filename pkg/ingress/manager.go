package ingress

import (
	"errors"

	"github.com/golang/glog"
	core_v1 "k8s.io/api/core/v1"
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
func NewDefaultManager() (Manager, error) {
	svc, err := NewLoadBalancerService()
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

	_, err := mgr.svc.CreateLoadBalancer(specification)
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

// getNodeInternalIPAddress will extract the OCI internal node IP address
// for a given node. Since it is impossible to launch an instance without
// an internal (private) IP, we can be sure that one exists.
func getNodeInternalIPAddress(node *core_v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == core_v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}
