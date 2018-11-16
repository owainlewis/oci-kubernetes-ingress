package controller

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"

	"github.com/owainlewis/oci-kubernetes-ingress/internal/config"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/context"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/store"

	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

// OCIController is the definition for an OCI Ingress Controller
type OCIController struct {
	configuration config.Config
	context       context.ControllerContext
	workQueue     OCIWorkQueue
	stopCh        chan struct{}

	nodeStore store.NodeStore
}

// NewOCIController will create a new OCI Ingress Controller
func NewOCIController(conf config.Config, context context.ControllerContext, stopCh chan struct{}) *OCIController {
	ctrl := &OCIController{
		configuration: conf,
		context:       context,
		stopCh:        stopCh,

		nodeStore: store.NewNodeStoreFromCache(context.CacheGroup.NodeCache),
	}

	ctrl.workQueue = NewOCIWorkQueue(ctrl.sync)

	// Ingress event handlers.
	ctrl.context.InformerGroup.IngressInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ingress := obj.(*extensions.Ingress)
			glog.V(4).Infof("Ingress %s added. Enqueing work item.", ingress.Name)
			ctrl.workQueue.Enqueue(ingress)
		},
		UpdateFunc: func(old, new interface{}) {
			newIngress := new.(*extensions.Ingress)
			ingressKeyName := fmt.Sprintf("%s/%s", newIngress.Namespace, newIngress.Name)

			if reflect.DeepEqual(old, new) {
				glog.V(4).Infof("Periodic enqueueing of %v", ingressKeyName)
			} else {
				glog.V(4).Infof("Ingress %s changed, enqueuing", ingressKeyName)
			}

			ctrl.workQueue.Enqueue(newIngress)
		},
		DeleteFunc: func(obj interface{}) {
			ingress := obj.(*extensions.Ingress)
			glog.V(4).Infof("Ingress %s deleted. Enqueing work item", ingress.Name)
			ctrl.workQueue.Enqueue(ingress)
		},
	})

	return ctrl
}

// Run will start the OCI Ingress Controller
func (c *OCIController) Run() {
	glog.Infof("Starting OCI Ingress controller")
	go c.context.Run()
	go c.workQueue.Run()

	<-c.stopCh
}

// sync manages Ingress create/updates/delete events from the work queue.
func (c *OCIController) sync(key string) error {
	if !c.context.HasSynced() {
		time.Sleep(5 * time.Second)
		return fmt.Errorf("waiting for cache stores to sync")
	}

	glog.V(4).Infof("\n\nSync called for item: %s\n\n", key)
	obj, exists, _ := c.context.CacheGroup.IngressCache.GetByKey(key)
	if !exists {
		glog.V(2).Infof("Ingress %v no longer exists. Garbage collecting OCI resources associated with it", key)
		return nil
	}

	// Get ingress and DeepCopy for assurance that we don't pollute other goroutines with changes.
	ingress, ok := obj.(*extensions.Ingress)
	if !ok {
		return fmt.Errorf("invalid object type %T. Expecting to be of type Ingress", obj)
	}
	ingress = ingress.DeepCopy()

	glog.V(4).Infof("Here is where we pass the ingress to a sync manager to ensure state in OCI %v", ingress)

	// We convert the Ingress object to a specification that descibes the desired state
	// of the OCI load balancer.
	// ingressSpecification := manager.NewSpecification(c.configuration, ingress, nil)

	return nil
}
