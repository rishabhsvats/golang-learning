package controller

import (
	"context"
	"log"
	"time"

	"github.com/rishabhsvats/golang-learning/kluster/pkg/apis/rishabhsvats.dev/v1alpha1"
	"github.com/rishabhsvats/golang-learning/kluster/pkg/do"
	klientset "github.com/rishabhsvats/golang-learning/kluster/pkg/generated/clientset/versioned"
	kinf "github.com/rishabhsvats/golang-learning/kluster/pkg/generated/informers/externalversions/rishabhsvats.dev/v1alpha1"
	klister "github.com/rishabhsvats/golang-learning/kluster/pkg/generated/listers/rishabhsvats.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	client kubernetes.Interface
	// clientset for custom resource kluster
	klient klientset.Interface
	// cluster cache has synced
	klusterSynced cache.InformerSynced
	// lister
	kLister klister.KlusterLister
	// queue
	wq workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, klient klientset.Interface, klusterInformer kinf.KlusterInformer) *Controller {
	c := &Controller{
		client:        client,
		klient:        klient,
		klusterSynced: klusterInformer.Informer().HasSynced,
		kLister:       klusterInformer.Lister(),
		wq:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "kluster"),
	}

	klusterInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	return c
}

func (c *Controller) Run(ch chan struct{}) error {
	if ok := cache.WaitForCacheSync(ch, c.klusterSynced); !ok {
		log.Println("cache was not synced")
	}

	go wait.Until(c.worker, time.Second, ch)

	<-ch
	return nil
}

func (c *Controller) worker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	item, shutdown := c.wq.Get()
	if shutdown {
		//log
		return false
	}
	defer c.wq.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		log.Printf("error %s calling namespace key func on cache item", err.Error())
		return false
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.Printf("splitting key into namespace and name, error %s\n", err.Error())
		return false
	}

	kluster, err := c.kLister.Klusters(ns).Get(name)
	if err != nil {
		log.Printf("error %s, getting the kluster resource from lister\n", err.Error())
		return false
	}

	log.Printf("kluster spec that we have is %+v\n", kluster.Spec)
	clusterID, err := do.Create(c.client, kluster.Spec)
	if err != nil {
		log.Printf("error %s, creating the cluster", err.Error())
	}
	log.Printf("cluster id that was created is %s\n", clusterID)

	err = c.updateStatus(clusterID, "creating", kluster)
	if err != nil {
		log.Printf("error  %s, updating status of the kluster %s\n", err.Error(), kluster.Name)
	}
	return true
}

func (c *Controller) updateStatus(id, progress string, kluster *v1alpha1.Kluster) error {
	kluster.Status.KlusterID = id
	kluster.Status.Progress = progress
	_, err := c.klient.RishabhsvatsV1alpha1().Klusters(kluster.Namespace).UpdateStatus(context.Background(), kluster, metav1.UpdateOptions{})
	return err
}

func (c *Controller) handleAdd(obj interface{}) {
	log.Println("handleAdd was called")
	c.wq.Add(obj)
}

func (c *Controller) handleDel(obj interface{}) {
	log.Println("handleDel was called")
	c.wq.Add(obj)
}
