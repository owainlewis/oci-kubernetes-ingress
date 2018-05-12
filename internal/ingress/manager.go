package ingress

import (
	"github.com/golang/glog"
	"k8s.io/api/extensions/v1beta1"
)

// Manager transforms ingress objects into OCI load balancer specifications
// and ensures that the desired and actual state of the world align.
type Manager interface {
	EnsureIngress(ingress *v1beta1.Ingress)
	EnsureIngressDeleted(ingress *v1beta1.Ingress)
}

type defaultManager struct {
}

// NewManager constructs a new ingress Manager
func NewManager() Manager {
	return &defaultManager{}
}

// EnsureIngress will ensure that the observed Ingress object state is reflected
// in OCI
func (manager *defaultManager) EnsureIngress(ingress *v1beta1.Ingress) {
	glog.Infof("Ensuring ingress")
	glog.Infof("Ingress Spec: %+v", ingress.Spec)
}

// EnsureIngressDeleted will ensure that an ingress object is removed from OCI
func (manager *defaultManager) EnsureIngressDeleted(ingress *v1beta1.Ingress) {
	glog.Infof("Deleting ingress")
}
