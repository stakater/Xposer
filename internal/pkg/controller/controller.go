package controller

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
	routeClient "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/sirupsen/logrus"
	config "github.com/stakater/Xposer/internal/pkg/config"
	"github.com/stakater/Xposer/internal/pkg/configmaps"
	"github.com/stakater/Xposer/internal/pkg/constants"
	"github.com/stakater/Xposer/internal/pkg/ingresses"
	"github.com/stakater/Xposer/internal/pkg/routes"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Event which is used to send additional information than just the key, can have other entiites
type Event struct {
	key       string
	eventType string
	oldObject interface{}
	newObject interface{}
}

// Controller for checking items
type Controller struct {
	clientset   kubernetes.Interface
	osClient    *routeClient.RouteV1Client
	clusterType string
	namespace   string
	indexer     cache.Indexer
	queue       workqueue.RateLimitingInterface
	informer    cache.Controller
	config      config.Configuration
}

// NewController A Constructor for the Controller to initialize the controller
func NewController(clientset kubernetes.Interface, osClient *routeClient.RouteV1Client, conf config.Configuration, clusterType string, namespace string) *Controller {
	controller := &Controller{
		clientset:   clientset,
		osClient:    osClient,
		config:      conf,
		clusterType: clusterType,
		namespace:   namespace,
	}

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	listWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), constants.SERVICES, namespace, fields.Everything())

	indexer, informer := cache.NewIndexerInformer(listWatcher, &v1.Service{}, constants.RESYNC_PERIOD, cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.Add,    //function that is called when the object is created
		UpdateFunc: controller.Update, //function that is called when the object is updated
		DeleteFunc: controller.Delete, //function that is called when the object is deleted
	}, cache.Indexers{})

	controller.indexer = indexer
	controller.informer = informer
	controller.queue = queue
	return controller
}

//Add function to add a 'create' event to the queue
func (c *Controller) Add(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	var event Event

	if err == nil {
		event.key = key
		event.eventType = "create"
		event.newObject = obj
		c.queue.Add(event)
	}
}

//Update function to add an 'update' event to the queue
func (c *Controller) Update(oldObj interface{}, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	var event Event

	if err == nil {
		event.key = key
		event.eventType = "update"
		event.oldObject = oldObj
		event.newObject = newObj
		c.queue.Add(event)
	}
}

//Delete function to add a 'delete' event to the queue
func (c *Controller) Delete(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	var event Event

	if err == nil {
		event.key = key
		event.eventType = "delete"
		event.newObject = obj
		c.queue.Add(event)
	}
}

//Run function for controller which handles the queue
func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	event, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two ingresses with the same key are never processed in
	// parallel.
	defer c.queue.Done(event)

	// Invoke the method containing the business logic
	err := c.takeAction(event.(Event))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, event)
	return true
}

//takeAction, the main function which will be handling the controller business logic
func (c *Controller) takeAction(event Event) error {
	// process events based on its type
	switch event.eventType {
	case "create":
		c.serviceCreated(event.newObject)

	case "update":
		c.serviceUpdated(event.oldObject, event.newObject)

	case "delete":
		c.serviceDeleted(event.newObject) //Incase of deleted, the obj object is nil
	}

	return nil
}

func (c *Controller) serviceCreated(obj interface{}) {
	newServiceObject := obj.(*v1.Service)
	logrus.Info("Service create event for the following service: %v", newServiceObject.Name)

	// Label for wether to create an ingress for this service or not
	if newServiceObject.ObjectMeta.Labels[constants.EXPOSE] == "true" {
		ingressInfo := ingresses.CreateIngressInfo(newServiceObject, c.config, c.namespace)

		if c.clusterType == constants.KUBERNETES {
			ingress := ingresses.CreateFromIngressInfo(ingressInfo)
			// Adds TLS for cert-manager if specified via annotations
			if ingressInfo.AddTLS == true {
				logrus.Info("Service contain TLS annotation, so automatically generating a TLS certificate via certmanager")
				ingresses.AddTLSInfo(ingress, ingressInfo.IngressName, ingressInfo.IngressHost)
			}

			result, err := c.clientset.ExtensionsV1beta1().Ingresses(c.namespace).Create(ingress)
			if err != nil {
				logrus.Warnf("Can not create new Ingress: %v", err)
			} else {
				logrus.Infof("Successfully created an Ingress with name: %v", result.Name)
			}

			if ingressInfo.ForwardAnnotationsMap[constants.EXPOSE_INGRESS_URL] == constants.GLOBALLY {
				configmaps.PopulateConfigMapGlobally(c.clientset, newServiceObject, ingressInfo.IngressHost)
			} else if ingressInfo.ForwardAnnotationsMap[constants.EXPOSE_INGRESS_URL] == constants.LOCALLY {
				configmaps.PopulateConfigMapLocally(c.clientset, newServiceObject, ingressInfo.IngressHost)
			}
		}

		if c.clusterType == constants.OPENSHIFT {
			route := routes.Create(ingressInfo.IngressName, ingressInfo.Namespace, ingressInfo.ForwardAnnotationsMap,
				ingressInfo.IngressHost, ingressInfo.IngressPath, ingressInfo.ServiceName, ingressInfo.ServicePort)

			result, err := c.osClient.Routes(c.namespace).Create(route)

			if err != nil {
				logrus.Errorf("Error while creating Route: %v", err)
			} else {
				logrus.Infof("Successfully created a Route with name: %v", result.Name)
			}
		}
	}
}

func (c *Controller) serviceUpdated(oldObj interface{}, newObj interface{}) {
	newServiceObject := newObj.(*v1.Service)
	oldServiceObject := oldObj.(*v1.Service)

	if oldServiceObject != newServiceObject {
		if newServiceObject.ObjectMeta.Labels[constants.EXPOSE] == "true" && oldServiceObject.ObjectMeta.Labels[constants.EXPOSE] == "true" {
			oldIngressConfig := structs.Map(c.config)
			oldIngressConfig = config.ReplaceDefaultConfigWithProvidedServiceConfig(oldIngressConfig, oldServiceObject)

			newIngressConfig := structs.Map(c.config)
			newIngressConfig = config.ReplaceDefaultConfigWithProvidedServiceConfig(newIngressConfig, newServiceObject)

			if oldIngressConfig[constants.INGRESS_NAME_TEMPLATE].(string) != newIngressConfig[constants.INGRESS_NAME_TEMPLATE].(string) {
				logrus.Info("Old service's Ingress Name template is different from new Service's Ingress Name Template. So deleting and re-creating Ingress in this case")
				c.serviceDeleted(oldObj)
				c.serviceCreated(newObj)
			} else {
				ingressInfo := ingresses.CreateIngressInfo(newServiceObject, c.config, c.namespace)
				ingress := ingresses.CreateFromIngressInfo(ingressInfo)

				if ingressInfo.AddTLS == true {
					ingresses.AddTLSInfo(ingress, ingressInfo.IngressName, ingressInfo.IngressHost)
					logrus.Info("Added TLS Info for certmanager")
				}

				result, err := c.clientset.ExtensionsV1beta1().Ingresses(c.namespace).Update(ingress)
				if err != nil {
					logrus.Errorf("Error while Updating Ingress: %v", err)
				} else {
					logrus.Infof("Successfully updated an Ingress with name: %v, for service: %v", result.Name, result.Spec.Backend.ServiceName)
				}

				// Updating exposed services URLs
				if ingressInfo.ForwardAnnotationsMap[constants.EXPOSE_INGRESS_URL] == constants.GLOBALLY {
					configmaps.DeleteFromConfigMapGlobally(c.clientset, oldServiceObject)
					configmaps.PopulateConfigMapGlobally(c.clientset, newServiceObject, ingressInfo.IngressHost)
				} else if ingressInfo.ForwardAnnotationsMap[constants.EXPOSE_INGRESS_URL] == constants.LOCALLY {
					configmaps.DeleteFromConfigMapLocally(c.clientset, oldServiceObject)
					configmaps.PopulateConfigMapLocally(c.clientset, newServiceObject, ingressInfo.IngressHost)
				}
			}
		} else {
			if newServiceObject.ObjectMeta.Labels[constants.EXPOSE] == "false" {
				c.serviceDeleted(oldObj)
			}

			if newServiceObject.ObjectMeta.Labels[constants.EXPOSE] == "true" {
				c.serviceCreated(newObj)
			}
		}
	} else {
		ingressList, err := c.clientset.ExtensionsV1beta1().Ingresses(c.namespace).List(meta_v1.ListOptions{})
		if err != nil {
			logrus.Errorf("Can not fetch Ingresses in the following namespace: %v, with the following error: %v", c.namespace, err)
		}
		existingIngress := ingresses.GetFromListMatchingGivenServiceName(ingressList, newServiceObject.Name)
		if ingresses.IsEmpty(existingIngress) {
			c.serviceCreated(newObj)
		}
	}
}

func (c *Controller) serviceDeleted(deletedServiceObject interface{}) {
	serviceToDelete := deletedServiceObject.(*v1.Service)
	logrus.Info("Service delete event for the following service: %v", serviceToDelete.Name)

	// Only delete ingress if the service had expose = true label
	if serviceToDelete.ObjectMeta.Labels["expose"] == "true" {

		ingressList, err := c.clientset.ExtensionsV1beta1().Ingresses(c.namespace).List(meta_v1.ListOptions{})
		if err != nil {
			logrus.Errorf("Can not fetch Ingresses in the following namespace: %v, with the following error: %v", c.namespace, err)
		}

		ingressToRemove := ingresses.GetFromListMatchingGivenServiceName(ingressList, serviceToDelete.Name)
		err = c.clientset.ExtensionsV1beta1().Ingresses(c.namespace).Delete(ingressToRemove.ObjectMeta.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			logrus.Warnf("Ingress not deleted with name: %v", ingressToRemove.ObjectMeta.Name)
		} else {
			logrus.Infof("Ingress Deleted with name: %v", ingressToRemove.ObjectMeta.Name)
		}

		// Updating xposer config map if it exists
		ingressInfo := ingresses.CreateIngressInfo(serviceToDelete, c.config, serviceToDelete.Namespace)

		if ingressInfo.ForwardAnnotationsMap[constants.EXPOSE_INGRESS_URL] == constants.GLOBALLY {
			configmaps.DeleteFromConfigMapGlobally(c.clientset, serviceToDelete)
		} else if ingressInfo.ForwardAnnotationsMap[constants.EXPOSE_INGRESS_URL] == constants.LOCALLY {

			configmaps.DeleteFromConfigMapLocally(c.clientset, serviceToDelete)
		}
	}
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		logrus.Printf("Error syncing ingress %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	logrus.Printf("Dropping ingress %q out of the queue: %v", key, err)
}
