package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/extensions/v1beta1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	lister_v1beta1 "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	kubeinformers "k8s.io/client-go/informers"

	"github.com/owainlewis/oci-ingress/internal/ingress"
)

type OCIController struct {
	client           kubernetes.Interface
	ingressLister    lister_v1beta1.IngressLister
	ingressWorkQueue workqueue.RateLimitingInterface
	ingressSynced    cache.InformerSynced

	ingressManager *ingress.Manager

	namespace string
}

func NewOCIController(client kubernetes.Interface, namespace string, informerFactory kubeinformers.SharedInformerFactory) *OCIController {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	ingressInformer := informerFactory.Extensions().V1beta1().Ingresses()

	ctrl := &OCIController{
		client:           client,
		ingressWorkQueue: queue,
		ingressLister:    ingressInformer.Lister(),
		ingressSynced:    ingressInformer.Informer().HasSynced,
		ingressManager:   ingress.NewManager(),
		namespace:        namespace,
	}

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.enqueueIngress,
		UpdateFunc: func(old, new interface{}) {
			newIngress := new.(*v1beta1.Ingress)
			oldIngress := old.(*v1beta1.Ingress)
			if newIngress.ResourceVersion == oldIngress.ResourceVersion {
				return
			}
			ctrl.enqueueIngress(new)
		},
		DeleteFunc: func(obj interface{}) {
			ingress, ok := obj.(*v1beta1.Ingress)
			if ok {
				ctrl.ingressManager.EnsureIngressDeleted(ingress)
			}
		},
	})

	return ctrl
}

func (c *OCIController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.ingressWorkQueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting OCI Ingress Controller")
	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.ingressSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *OCIController) enqueueIngress(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}

	c.ingressWorkQueue.AddRateLimited(key)
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *OCIController) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *OCIController) processNextWorkItem() bool {
	obj, shutdown := c.ingressWorkQueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.ingressWorkQueue.Done(obj)
		key, ok := obj.(string)
		if !ok {
			c.ingressWorkQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in queue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.ingressWorkQueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *OCIController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	ingress, err := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil {
		return err
	}

	c.ingressManager.EnsureIngress(ingress)
	return nil
}
