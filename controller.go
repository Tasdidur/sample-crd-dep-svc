package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	//v1 "k8s.io/client-go/informers/core/v1"
	//"k8s.io/client-go/listers/extensions/v1beta1"

	//"github.com/appscode/go/runtime"
	//"github.com/appscode/go/wait"
	xapiv1 "github.com/Tasdidur/xcrd/pkg/apis/xapi.com/v1"
	clientset "github.com/Tasdidur/xcrd/pkg/client/clientset/versioned"
	informers "github.com/Tasdidur/xcrd/pkg/client/informers/externalversions/xapi.com/v1"
	lister "github.com/Tasdidur/xcrd/pkg/client/listers/xapi.com/v1"
	//appv1 "k8s.io/api/apps/v1"
	//corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	corev1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	appslister "k8s.io/client-go/listers/apps/v1"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
)

const controllerAgentName = "CRD-Controller"

const (
	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	MessageResourceSynced = "Foo synced successfully"
)

// Controller is the controller implementation for Foo resources
type Controller struct {
	kubeclientset   kubernetes.Interface
	sampleclientset clientset.Interface

	deploymentsLister appslister.DeploymentLister
	deploymentsSynced cache.InformerSynced

	serviceLister corev1lister.ServiceLister
	serviceSynced cache.InformerSynced

	xcrdLister        lister.XcrdLister
	xcrdSynced        cache.InformerSynced

	workQueue workqueue.RateLimitingInterface
}

// NewController returns a new sample controller

func NewController(kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformer.DeploymentInformer,
	serviceInformer corev1informer.ServiceInformer,
	xcrdInformer informers.XcrdInformer) *Controller {

	controller := &Controller{
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		deploymentsLister: deploymentInformer.Lister(),
		serviceSynced: 	   serviceInformer.Informer().HasSynced,
		serviceLister: 	   serviceInformer.Lister(),
		xcrdLister:        xcrdInformer.Lister(),
		xcrdSynced:        xcrdInformer.Informer().HasSynced,
		workQueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Xcrds"),
	}

	log.Println("setting up event handlers")

	xcrdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueXcrd,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueXcrd(new)
		},
	})
	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workQueue.ShutDown()

	log.Println("Starting xcrd Controller")

	log.Println("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.xcrdSynced); !ok {
		return fmt.Errorf("failed to wait for cache to sync")

	}
	log.Println("Starting workers")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Println("Worker started")
	<-stopCh
	log.Println("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.ProcessNextItem() {

	}

}

func (c *Controller) ProcessNextItem() bool {
	obj, shutdown := c.workQueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workQueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workQueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil

		}
		if err := c.syncHandler(key); err != nil {
			c.workQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		c.workQueue.Forget(obj)
		log.Printf("successfully synced '%s'\n", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}
	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
	}

	xcrd, err := c.xcrdLister.Xcrds(namespace).Get(name)

	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("xcrd '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	deploymentName := xcrd.Spec.Name+"-dep"

	deployment, err := c.deploymentsLister.Deployments(namespace).Get(deploymentName)

	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(xcrd.Namespace).Create(context.TODO(),newDeployment(xcrd),metav1.CreateOptions{})
	}

	fmt.Println(deployment.Name)

	if err != nil {
		return err
	}

	serviceName := xcrd.Spec.Name+"-svc"

	service , err := c.serviceLister.Services(namespace).Get(serviceName)

	if errors.IsNotFound(err) {
		service , err = c.kubeclientset.CoreV1().Services(xcrd.Namespace).Create(context.TODO(),newService(xcrd),metav1.CreateOptions{})
	}

	fmt.Println(service.Name)

	if err != nil {
		return err
	}

	err = c.updateXcrdStatus(xcrd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) enqueueXcrd(obj interface{}) {
	log.Println("Enqueueing Xcrd. . . ")
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workQueue.AddRateLimited(key)
}

func (c *Controller) updateXcrdStatus(xcrd *xapiv1.Xcrd) error {
	// call if all available

	xcrdCopy := xcrd.DeepCopy()
	xcrdCopy.Status.AllReady = true
	_, err := c.sampleclientset.XapiV1().Xcrds(xcrd.Namespace).Update(context.TODO(),xcrdCopy,metav1.UpdateOptions{})

	return err
}

//func newDeployment(foo *controllerv1alpha1.Foo) *appv1.Deployment {
//
//	return &appv1.Deployment{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      foo.Spec.DeploymentName,
//			Namespace: foo.Namespace,
//		},
//		Spec: appv1.DeploymentSpec{
//			Replicas: foo.Spec.Replicas,
//			Selector: &metav1.LabelSelector{
//				MatchLabels: map[string]string{
//					"app":        "nginx",
//					"controller": foo.Name,
//				},
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels: map[string]string{
//						"app":        "nginx",
//						"controller": foo.Name,
//					},
//				},
//				Spec: corev1.PodSpec{
//					Containers: []corev1.Container{
//						{
//							Name:  "nginx",
//							Image: "nginx:latest",
//						},
//					},
//				},
//			},
//		},
//	}
//}
