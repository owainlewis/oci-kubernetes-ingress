package ingress

import (
	"github.com/golang/glog"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

// Manager transforms ingress objects into OCI load balancer specifications
// and ensures that the desired and actual state of the world align.
type Manager interface {
	EnsureIngress(ingress *v1beta1.Ingress, nodes []*core_v1.Node)
	EnsureIngressDeleted(ingress *v1beta1.Ingress)
}

type defaultManager struct {
}

// NewDefaultManager constructs a new ingress Manager
func NewDefaultManager() Manager {
	return &defaultManager{}
}

// EnsureIngress will ensure that the observed Ingress object state is reflected
// in OCI
func (manager *defaultManager) EnsureIngress(ingress *v1beta1.Ingress, nodes []*core_v1.Node) {
	glog.Infof("Ensuring ingress")
	glog.Infof("Ingress Spec: %+v", ingress.Spec)
}

// EnsureIngressDeleted will ensure that an ingress object is removed from OCI
func (manager *defaultManager) EnsureIngressDeleted(ingress *v1beta1.Ingress) {
	glog.Infof("Deleting ingress")
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
