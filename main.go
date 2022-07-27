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




func main() {
}






func (c *controller) Run(stopCh <-chan struct{}) {
}





func (c *controller) runWorker() {
}





func (c *controller) processNextItem() bool {
}





func (c *controller) processItem(key string) error {
}










type controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	client   clientset.Interface

	craiyon *craiyon.Client
	// deepai *deepai.Client
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
