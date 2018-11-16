package context

import (
	"fmt"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InformerGroup is a structure that holds multiple informers.
type InformerGroup struct {
	IngressInformer cache.SharedIndexInformer
	ServiceInformer cache.SharedIndexInformer
	SecretInformer  cache.SharedIndexInformer
}

// Run will run and await sync for all the informers in the informer group.
func (i *InformerGroup) Run(stopCh chan struct{}) error {
	go i.ServiceInformer.Run(stopCh)
	go i.SecretInformer.Run(stopCh)
	go i.IngressInformer.Run(stopCh)

	glog.V(4).Info("Waiting for all caches to sync")

	if !cache.WaitForCacheSync(stopCh, i.IngressInformer.HasSynced, i.ServiceInformer.HasSynced, i.SecretInformer.HasSynced) {
		return fmt.Errorf("failed waiting for caches to sync")
	}

	glog.V(4).Info("All caches have synced")

	return nil
}

// CacheGroup is a structure that holds multiple caches.
type CacheGroup struct {
	IngressCache cache.Store
	ServiceCache cache.Store
	SecretCache  cache.Store
}

// ControllerContext provides a controller with access to multiple informers and caches.
type ControllerContext struct {
	InformerGroup InformerGroup
	CacheGroup    CacheGroup
	StopChannel   chan struct{}
}

// NewControllerContext will construct a new ControllerContext struct.
func NewControllerContext(kubeClient kubernetes.Interface, namespace string, resyncPeriod time.Duration) ControllerContext {
	informerFactory := informers.NewFilteredSharedInformerFactory(kubeClient, resyncPeriod, namespace, func(*metav1.ListOptions) {})

	ctx := ControllerContext{
		InformerGroup: InformerGroup{
			IngressInformer: informerFactory.Extensions().V1beta1().Ingresses().Informer(),
			ServiceInformer: informerFactory.Core().V1().Services().Informer(),
			SecretInformer:  informerFactory.Core().V1().Secrets().Informer(),
		},
		CacheGroup:  CacheGroup{},
		StopChannel: make(chan struct{}),
	}

	ctx.CacheGroup.IngressCache = ctx.InformerGroup.IngressInformer.GetStore()
	ctx.CacheGroup.ServiceCache = ctx.InformerGroup.ServiceInformer.GetStore()
	ctx.CacheGroup.SecretCache = ctx.InformerGroup.SecretInformer.GetStore()

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
		c.InformerGroup.SecretInformer.HasSynced() &&
		c.InformerGroup.ServiceInformer.HasSynced()
}
