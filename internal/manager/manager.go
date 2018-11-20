package manager

import (
	"github.com/golang/glog"
)

// IngressManager is an interface for managing state between K8s Ingress and OCI LBs.
type IngressManager interface {
	// EnsureIngress will ensure that an OCI load balancer exists and is
	// configured correctly for the provided ingress object.
	EnsureIngress()
	// EnsureIngressDeleted ensures that all OCI resources associated with
	// an Ingress are removed.
	EnsureIngressDeleted()
}

type OCIIngressManager struct {
}

func NewOCIIngressManager() *OCIIngressManager {
	return &OCIIngressManager{}
}

func (m *OCIIngressManager) EnsureIngress() {
	glog.V(4).Info("Ensuring ingress state")

}

func (m *OCIIngressManager) EnsureIngressDeleted() {
	glog.V(4).Info("Ensuring ingress deleted")
}
