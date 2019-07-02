package handlers

import (
	"context"

	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/ingress/annotations"
)

// When nodes change (added or removed from cluster) update the state of the load balancers for all ingress objects

var _ handler.EventHandler = (*EnqueueRequestsForNodeEvent)(nil)

type EnqueueRequestsForNodeEvent struct {
	Cache cache.Cache
}

func (h *EnqueueRequestsForNodeEvent) Create(e event.CreateEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueImpactedIngresses(queue)
}

func (h *EnqueueRequestsForNodeEvent) Delete(e event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueImpactedIngresses(queue)
}

func (h *EnqueueRequestsForNodeEvent) Update(e event.UpdateEvent, queue workqueue.RateLimitingInterface) {

}

func (h *EnqueueRequestsForNodeEvent) Generic(event.GenericEvent, workqueue.RateLimitingInterface) {
}

func (h *EnqueueRequestsForNodeEvent) enqueueImpactedIngresses(queue workqueue.RateLimitingInterface) {
	ingressList := &extensions.IngressList{}
	if err := h.Cache.List(context.Background(), ingressList); err != nil {
		return
	}

	for _, ingress := range ingressList.Items {
		if annotations.HasOCIIngressAnnotation(&ingress) {
			queue.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: ingress.Namespace,
					Name:      ingress.Name,
				},
			})
		}
	}
}
