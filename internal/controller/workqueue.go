package controller

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/golang/glog"
)

// OCIWorkQueue is a struct representing a controller work queue.
type OCIWorkQueue struct {
	queue workqueue.RateLimitingInterface
	// processItem is a function called for each item in the work queue.
	sync func(string) error
	// workerDone channel is closed when the queue worker exits.
	workerDone chan struct{}
}

// NewOCIWorkQueue constructs a new OCIWorkQueue.
func NewOCIWorkQueue(sync func(string) error) OCIWorkQueue {
	return OCIWorkQueue{
		queue:      workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		sync:       sync,
		workerDone: make(chan struct{}),
	}
}

// Enqueue adds an object to the work queue.
func (q OCIWorkQueue) Enqueue(objs ...interface{}) {
	glog.V(4).Infof("Queue depth at %d", q.queue.Len())

	for _, obj := range objs {
		var key string
		var err error

		if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
			glog.Errorf("Couldn't get cache key for object: %v", obj)
			return
		}

		glog.V(4).Infof("Enqueuing object with key: %v", key)
		q.queue.Add(key)
		glog.V(4).Infof("Queue depth at %d", q.queue.Len())
	}
}

// Run will continue to process items from the work queue.
func (q OCIWorkQueue) Run() {
	for {
		key, quit := q.queue.Get()
		if quit {
			close(q.workerDone)
			return
		}
		glog.V(4).Infof("Syncing %v", key)
		if err := q.sync(key.(string)); err != nil {
			glog.Errorf("Requeuing %q due to error", key)
			q.queue.AddRateLimited(key)
		} else {
			glog.V(4).Infof("Finished syncing %v", key)
			q.queue.Forget(key)
		}
		q.queue.Done(key)
	}
}

// ShutDown will shutdown the work queue.
func (q OCIWorkQueue) ShutDown() {
	glog.V(2).Infof("Shutdown called on work queue")
	q.queue.ShutDown()
	<-q.workerDone
}
