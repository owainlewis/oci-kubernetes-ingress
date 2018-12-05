package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/client-go/listers/core/v1"
	v1beta1 "k8s.io/client-go/listers/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/labels"

	apiv1 "k8s.io/api/core/v1"
)

// InformerGroup is a structure that holds multiple informers.
type InformerGroup struct {
	IngressInformer cache.SharedIndexInformer
	NodeInformer    cache.SharedIndexInformer
	ServiceInformer cache.SharedIndexInformer
}

// Run will run and await sync for all the informers in the informer group.
func (i *InformerGroup) Run(stopCh chan struct{}) error {
	go i.IngressInformer.Run(stopCh)
	go i.NodeInformer.Run(stopCh)
	go i.ServiceInformer.Run(stopCh)

	glog.V(4).Info("Waiting for all caches to sync")

	if !cache.WaitForCacheSync(stopCh,
		i.IngressInformer.HasSynced,
		i.NodeInformer.HasSynced,
		i.ServiceInformer.HasSynced) {
		return fmt.Errorf("failed waiting for caches to sync")
	}

	glog.V(4).Info("All caches have synced")

	return nil
}

// ListerGroup is a structure that holds multiple listers.
type ListerGroup struct {
	IngressLister v1beta1.IngressLister
	NodeLister    corev1.NodeLister
	ServiceLister corev1.ServiceLister
}

// CacheGroup is a structure that holds multiple caches.
type CacheGroup struct {
	IngressCache cache.Store
	NodeCache    cache.Store
	ServiceCache cache.Store
}

// ControllerContext provides a controller with access to multiple informers and caches.
type ControllerContext struct {
	InformerGroup InformerGroup
	ListerGroup   ListerGroup
	CacheGroup    CacheGroup
	StopChannel   chan struct{}
}

// NewControllerContext will construct a new ControllerContext struct.
func NewControllerContext(kubeClient kubernetes.Interface, namespace string, resyncPeriod time.Duration) ControllerContext {
	informerFactory := informers.NewFilteredSharedInformerFactory(kubeClient, resyncPeriod, namespace, func(*metav1.ListOptions) {})
	ctx := ControllerContext{
		InformerGroup: InformerGroup{
			IngressInformer: informerFactory.Extensions().V1beta1().Ingresses().Informer(),
			NodeInformer:    informerFactory.Core().V1().Nodes().Informer(),
			ServiceInformer: informerFactory.Core().V1().Services().Informer(),
		},
		ListerGroup: ListerGroup{
			IngressLister: informerFactory.Extensions().V1beta1().Ingresses().Lister(),
			NodeLister:    informerFactory.Core().V1().Nodes().Lister(),
			ServiceLister: informerFactory.Core().V1().Services().Lister(),
		},
		StopChannel: make(chan struct{}),
	}

	ctx.CacheGroup.IngressCache = ctx.InformerGroup.IngressInformer.GetStore()
	ctx.CacheGroup.NodeCache = ctx.InformerGroup.NodeInformer.GetStore()
	ctx.CacheGroup.ServiceCache = ctx.InformerGroup.ServiceInformer.GetStore()

	return ctx
}

// Run will run all informers in the informer context.
func (c *ControllerContext) Run() {
	glog.V(4).Info("Running informers")
	c.InformerGroup.Run(c.StopChannel)
}

// Stop will stop all the informers in the informer group.
func (c *ControllerContext) Stop() {
	c.StopChannel <- struct{}{}
}

// HasSynced returns true if all informers have synced.
func (c *ControllerContext) HasSynced() bool {
	return c.InformerGroup.IngressInformer.HasSynced() &&
		c.InformerGroup.NodeInformer.HasSynced() &&
		c.InformerGroup.ServiceInformer.HasSynced()
}

// GetAllNodes returns a list of all nodes in the current namespace.
func (c *ControllerContext) GetAllNodes() ([]*apiv1.Node, error) {
	return c.ListerGroup.NodeLister.List(labels.Everything())
}

// GetAllServices returns a list of all services in the current namespace.
func (c *ControllerContext) GetAllServices() ([]*apiv1.Service, error) {
	return c.ListerGroup.ServiceLister.List(labels.Everything())
}
