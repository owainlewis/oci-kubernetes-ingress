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

type EnqueueRequestsForIngressEvent struct {
}

func (h *EnqueueRequestsForIngressEvent) Create(e event.CreateEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.Object.(*extensions.Ingress), queue)
}

func (h *EnqueueRequestsForIngressEvent) Update(e event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.ObjectOld.(*extensions.Ingress), queue)
	h.enqueueIfIngressClassMatched(e.ObjectNew.(*extensions.Ingress), queue)
}

func (h *EnqueueRequestsForIngressEvent) Delete(e event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	h.enqueueIfIngressClassMatched(e.Object.(*extensions.Ingress), queue)
}

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
