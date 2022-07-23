package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	picturesv1 "github.com/eddiezane/that-conference-k8s-controller/pkg/apis/pictures/v1"
	clientset "github.com/eddiezane/that-conference-k8s-controller/pkg/generated/clientset/versioned"
	informers "github.com/eddiezane/that-conference-k8s-controller/pkg/generated/informers/externalversions"

	craiyon "github.com/eddiezane/that-conference-k8s-controller/pkg/craiyon"
	// deepai "github.com/eddiezane/that-conference-k8s-controller/pkg/text2image"
)

type controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	client   clientset.Interface

	craiyon *craiyon.Client
	// deepai *deepai.Client
}

func (c *controller) Run(stopCh <-chan struct{}) {
	go c.informer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for sync"))
		return
	}

	klog.Info("starting worker")
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *controller) processNextItem() bool {
	key, shutdown := c.queue.Get()
	if shutdown {
		klog.InfoS("shutting down worker")
		return false
	}
	go func(key string) {
		defer c.queue.Done(key)

		err := c.processItem(key)
		if err == nil {
			klog.InfoS("done with key", "key", key)
			c.queue.Forget(key)
		} else if c.queue.NumRequeues(key) < 3 {
			utilruntime.HandleError(err)
			klog.InfoS("requeuing", "key", key)
			c.queue.AddRateLimited(key)
		} else {
			utilruntime.HandleError(err)
			klog.InfoS("done trying", "key", key)
			c.queue.Forget(key)
		}
	}(key.(string))

	return true
}

func (c *controller) processItem(key string) error {
	item, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("unable to fetch key %s from store %w", key, err)
	}

	if !exists {
		klog.InfoS("item deleted", "key", key)
	}

	picture := item.(*picturesv1.Picture)
	if err != nil {
		return fmt.Errorf("unable to split meta: %w", err)
	}

	if picture.Status.URL == "" || picture.Generation != picture.Status.ObservedGeneration {
		klog.InfoS("need to sync picture", "key", key)

		// url, err := c.deepai.GetImage(picture.Spec.Text)
		url, err := c.craiyon.GetImage(picture.Spec.Text)

		if err != nil {
			return fmt.Errorf("unable to get image url: %w", err)
		}
		p := picture.DeepCopy()
		p.Status.URL = url
		p.Status.ObservedGeneration = p.Generation
		_, err = c.client.KuberneddiesV1().Pictures(p.Namespace).UpdateStatus(context.Background(), p, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("unable to update picture status: %w", err)
		}
	} else {
		klog.InfoS("nothing to do", "key", key)
	}

	return nil
}

func main() {
	klog.InitFlags(nil)

	flags := genericclioptions.NewConfigFlags(true)
	config, err := flags.ToRESTConfig()
	if err != nil {
		panic(err)
	}

	client := clientset.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactory(client, 30*time.Second)
	informer := factory.Kuberneddies().V1().Pictures().Informer()

	c := NewController(
		workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		informer,
		client,
	)

	informer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.enqueue(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			old := oldObj.(*picturesv1.Picture)
			new := newObj.(*picturesv1.Picture)

			if old.ResourceVersion == new.ResourceVersion {
				klog.InfoS("periodic resync", "namespace", new.Namespace, "name", new.Name)
				return
			}
			c.enqueue(new)
		},
		DeleteFunc: func(obj interface{}) {
			p := obj.(*picturesv1.Picture)
			klog.InfoS("resource deleted", "namespace", p.Namespace, "name", p.Name)
		},
	})

	signalHandler := signals.SetupSignalHandler()

	klog.Info("starting controller")

	c.Run(signalHandler.Done())

	klog.Info("stopping controller")
}

func NewController(
	queue workqueue.RateLimitingInterface,
	informer cache.SharedIndexInformer,
	client clientset.Interface,
) *controller {
	return &controller{
		queue:    queue,
		informer: informer,
		client:   client,
		craiyon:  craiyon.NewClient(),
		// deepai:   deepai.NewClient(),
	}
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	klog.InfoS("adding to queue", "key", key)
	c.queue.Add(key)
}
