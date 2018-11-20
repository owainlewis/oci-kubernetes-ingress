package manager

import (
	"github.com/golang/glog"
	extensions "k8s.io/api/extensions/v1beta1"

	corev1 "k8s.io/api/core/v1"
)

// IngressManager is an interface for managing state between K8s Ingress and OCI LBs.
type IngressManager interface {
	// EnsureIngress will ensure that an OCI load balancer exists and is
	// configured correctly for the provided ingress object.
	EnsureIngress(ingress *extensions.Ingress, nodes []*corev1.Node)
	// EnsureIngressDeleted ensures that all OCI resources associated with
	// an Ingress are removed.
	EnsureIngressDeleted()
}

type OCIIngressManager struct {
}

func NewOCIIngressManager() *OCIIngressManager {
	return &OCIIngressManager{}
}

func (m *OCIIngressManager) EnsureIngress(ingress *extensions.Ingress, nodes []*corev1.Node) {
	glog.V(4).Infof("Ensuring state for ingress: %v for cluster with nodes %v", ingress, nodes)
}

func (m *OCIIngressManager) EnsureIngressDeleted() {
	glog.V(4).Info("Ensuring ingress deleted")
}
