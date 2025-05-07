package main

import (
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset     kubernetes.Interface
	depLister     appslisters.DeploymentLister
	depCacheSyncd cache.InformerSynced
	queue         workqueue.RateLimitingInterface
}

func newController() *controller {
	c := &controller{}
	return c
}
