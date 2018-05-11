package ingress

import (
	"github.com/golang/glog"
	"k8s.io/api/extensions/v1beta1"
)

// Manager maps desired to actual state for ingress objects in
// Kubernetes and Load Balancers in OCI
type Manager struct {
}

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
