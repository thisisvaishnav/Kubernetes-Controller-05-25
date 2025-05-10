package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	clientset *kubernetes.Clientset
	informer  cache.SharedIndexInformer
	workqueue workqueue.RateLimitingInterface
}

func NewController(clientset *kubernetes.Clientset, informer cache.SharedIndexInformer) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return &Controller{
		clientset: clientset,
		informer:  informer,
		workqueue: queue,
	}
}

func (c *Controller) Run(stopCh chan struct{}) {
	defer c.workqueue.ShutDown()

	fmt.Println("Starting Controller...")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		log.Fatalf("Failed to sync cache")
	}

	for {
		key, shutdown := c.workqueue.Get()
		if shutdown {
			break
		}
		fmt.Printf("Processing: %s\n", key)

		retry.RetryOnConflict(retry.DefaultRetry, func() error {
			obj, exists, err := c.informer.GetIndexer().GetByKey(key.(string))
			if err != nil {
				log.Printf("Error fetching object from store: %v\n", err)
				return err
			}
			if exists {
				pod := obj.(*v1.Pod)
				fmt.Printf("Pod Found: %s in namespace %s\n", pod.Name, pod.Namespace)
			} else {
				fmt.Printf("Pod Deleted: %s\n", key)
			}
			return nil
		})
		c.workqueue.Done(key)
	}
}

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	informer := factory.Core().V1().Pods().Informer()

	controller := NewController(clientset, informer)

	stopCh := make(chan struct{})
	defer close(stopCh)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	go controller.Run(stopCh)

	<-signalChan
	fmt.Println("Shutdown signal received, exiting...")
}
