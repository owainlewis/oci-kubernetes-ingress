package controller

import (
	"fmt"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/golang/glog"
)

// OCIWorkQueue is a struct representing a controller work queue.
type OCIWorkQueue struct {
	queue workqueue.RateLimitingInterface
	// processItem is a function called for each item in the work queue.
	processItem func(string) error
	// workerDone channel is closed when the queue worker exits.
	workerDone chan struct{}
}

// NewOCIWorkQueue constructs a new OCIWorkQueue.
func NewOCIWorkQueue(handler func(string) error) OCIWorkQueue {
	return OCIWorkQueue{
		queue:       workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		processItem: handler,
		workerDone:  make(chan struct{}),
	}
}

// Enqueue adds an object to the work queue.
func (q OCIWorkQueue) Enqueue(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}

	glog.V(4).Infof("Enqueuing object: %v", obj)
	q.queue.AddRateLimited(key)
}

// Run will continue to process items from the work queue.
func (q OCIWorkQueue) Run() {
	for {
		obj, shutdown := q.queue.Get()
		if shutdown {
			glog.V(4).Info("Queue is done")
			close(q.workerDone)
			return
		}

		defer q.queue.Done(obj)

		key, ok := obj.(string)

		glog.V(4).Infof("\n\n\nProcessing item %v\n\n\n", key)

		if !ok {
			q.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in queue but got %#v", obj))
		}
		if err := q.processItem(key); err != nil {
			glog.Errorf("Requeuing %q due to error: %v", key, err)
			q.queue.AddRateLimited(key)
		} else {
			glog.V(4).Infof("Successfully processed item: '%s'", key)
			q.queue.Forget(key)
		}

		q.queue.Forget(obj)
	}
}

// ShutDown will shutdown the work queue.
func (q OCIWorkQueue) ShutDown() {
	glog.V(2).Infof("Shutdown called on work queue")
	q.queue.ShutDown()
	<-q.workerDone
}
