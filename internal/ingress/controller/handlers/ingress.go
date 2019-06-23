package handlers

import (
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ handler.EventHandler = (*EnqueueRequestsForIngressEvent)(nil)

// EnqueueRequestsForIngressEvent defined a structure for ingress events
type EnqueueRequestsForIngressEvent struct {
}

// Create is a handler called when an integress object is created
func (h *EnqueueRequestsForIngressEvent) Create(e event.CreateEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.Object.(*extensions.Ingress), queue)
}

// Update is a handler called when an integress object is updated
func (h *EnqueueRequestsForIngressEvent) Update(e event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.ObjectOld.(*extensions.Ingress), queue)
	h.enqueueIfIngressClassMatched(e.ObjectNew.(*extensions.Ingress), queue)
}

// Delete is a handler called when an integress object is deleted
func (h *EnqueueRequestsForIngressEvent) Delete(e event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.Object.(*extensions.Ingress), queue)
}

// Generic is a handler called when an integress object is modified but none of the above
func (h *EnqueueRequestsForIngressEvent) Generic(e event.GenericEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.Object.(*extensions.Ingress), queue)
}

func (h *EnqueueRequestsForIngressEvent) enqueueIfIngressClassMatched(ingress *extensions.Ingress, queue workqueue.RateLimitingInterface) {
	// TODO only add to queue if the ingress class is appropriate
	queue.Add(reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: ingress.Namespace,
			Name:      ingress.Name,
		},
	})
}
