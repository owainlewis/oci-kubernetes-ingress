package ingress

import (
	"github.com/golang/glog"
	"k8s.io/api/extensions/v1beta1"
)

// Manager transforms ingress objects into OCI load balancer specifications
// and ensures that the desired and actual state of the world align.
type Manager struct {
}

// NewManager constructs a new ingress Manager
func NewManager() *Manager {
	return &Manager{}
}

// EnsureIngress will ensure that the observed Ingress object state is reflected
// in OCI
func (manager *Manager) EnsureIngress(ingress *v1beta1.Ingress) {
	glog.Infof("Ensuring ingress")
	glog.Infof("Ingress Spec: %+v", ingress.Spec)
}

// EnsureIngressDeleted will ensure that an ingress object is removed from OCI
func (manager *Manager) EnsureIngressDeleted(ingress *v1beta1.Ingress) {
	glog.Infof("Deleting ingress")
}
