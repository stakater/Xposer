package controller

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/stakater/Xposer/pkg/config"
	"github.com/stakater/Xposer/pkg/constants"
	"github.com/stakater/Xposer/pkg/ingresses"
	"github.com/stakater/Xposer/pkg/templates"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// AllNamespaces as our controller will be looking for events in all namespaces
const (
	AllNamespaces = ""
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
	clientset kubernetes.Interface
	resource  string
	namespace string
	indexer   cache.Indexer
	queue     workqueue.RateLimitingInterface
	informer  cache.Controller
	config    config.Configuration
}

// NewController A Constructor for the Controller to initialize the controller
func NewController(clientset kubernetes.Interface, conf config.Configuration, resource string, namespace string) *Controller {
	controller := &Controller{
		clientset: clientset,
		config:    conf,
		resource:  resource,
		namespace: namespace,
	}

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	// replace 'first-controller' with generic namespace after testing
	listWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), resource, "first-controller", fields.Everything())

	indexer, informer := cache.NewIndexerInformer(listWatcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
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

	fmt.Println("Adding create event to queue")
	fmt.Println("Service Name: ", obj.(*v1.Service).Name)
	fmt.Println("Service Spec: ", obj.(*v1.Service).Spec.Selector["app"])

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

	fmt.Println("Adding update event to queue")

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

	fmt.Println("Adding delete event to queue")

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
	//Printing Pod Name and its Containers from all namespaces but we can do anything in these functions
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
	// Currently printing all pods but we can restrict using any of the data in yaml file
	// e.g If want to check on APIVersion uncomment this.
	// if obj.(*v1.Pod).APIVersion == "samplecontroller.k8s.io/v1alpha1" {

	// fmt.Println("\nService Created Event")
	// fmt.Println("Name: ", obj.(*v1.Service).Name)
	// fmt.Println("First Annotation: ", obj.(*v1.Service).ObjectMeta.Annotations["firstAnnotation"])
	// fmt.Println("Second Annotation: ", strings.Split(obj.(*v1.Service).ObjectMeta.Annotations["xposer.stakater.com/annotations"], "\n"))
	// fmt.Println("Label 1: ", obj.(*v1.Service).ObjectMeta.Labels["k8sapp"])
	// fmt.Println("Label Expose: ", obj.(*v1.Service).ObjectMeta.Labels["expose"])
	// fmt.Println("Selector: ", obj.(*v1.Service).Spec.Ports[0].Port)
	// fmt.Println("Printing annotations")
	newServiceObject := obj.(*v1.Service)

	// Label for wether to create an ingress for this service or not
	if newServiceObject.ObjectMeta.Labels["expose"] == "true" {
		splittedAnnotations := strings.Split(string(newServiceObject.ObjectMeta.Annotations[constants.FORWARD_ANNOTATIONS]), "\n")
		fmt.Println("Splitted Annotations length: ", len(splittedAnnotations))

		forwardAnnotationsMap := make(map[string]string)
		ingressConfig := structs.Map(c.config)

		for annotationKey, annotationValue := range newServiceObject.ObjectMeta.Annotations {
			if strings.HasPrefix(annotationKey, constants.INGRESS_ANNOTATIONS) {
				ingressConfig[strings.SplitN(annotationKey, "/", 2)[1]] = annotationValue
			}
		}

		// Adds "/" in URL Path, if user has entered path annotaion without "/"
		if !strings.HasPrefix(ingressConfig[constants.INGRESS_URL_PATH].(string), "/") {
			ingressConfig[constants.INGRESS_URL_PATH] = "/" + ingressConfig[constants.INGRESS_URL_PATH].(string)
		}

		/*
			Removes the content after "/" from URL-Template, and if user has not specified path from annotation,
			use the content after "/" as URL-Path
		*/
		if strings.Contains(ingressConfig[constants.INGRESS_URL_TEMPLATE].(string), "/") {
			fmt.Println("Contains slash")
			splittedURLTemplate := strings.SplitN(ingressConfig[constants.INGRESS_URL_TEMPLATE].(string), "/", 2)
			ingressConfig[constants.INGRESS_URL_TEMPLATE] = splittedURLTemplate[0]

			if ingressConfig[constants.INGRESS_URL_PATH].(string) == "/" {
				ingressConfig[constants.INGRESS_URL_PATH] = ingressConfig[constants.INGRESS_URL_PATH].(string) + splittedURLTemplate[1]
			}
		}

		for _, annotation := range splittedAnnotations {
			fmt.Println("Annotation-split: ", annotation)
			parsedAnnotation := strings.Split(annotation, ":")
			if len(parsedAnnotation) != 2 {
				// throw error
			}
			forwardAnnotationsMap[parsedAnnotation[0]] = parsedAnnotation[1]
		}

		urlTemplate := &templates.URLTemplate{
			Service:   newServiceObject.Name,
			Namespace: "first-controller",
			Domain:    ingressConfig["Domain"].(string),
		}

		nameTemplate := &templates.NameTemplate{
			Service:   newServiceObject.Name,
			Namespace: "first-controller",
		}

		parsedURL := templates.ParseIngressURLOrPathTemplate(ingressConfig[constants.INGRESS_URL_TEMPLATE].(string), urlTemplate)
		parsedURLPath := templates.ParseIngressURLOrPathTemplate(ingressConfig[constants.INGRESS_URL_PATH].(string), urlTemplate)
		parsedIngressName := templates.ParseIngressNameTemplate(ingressConfig[constants.INGRESS_NAME_TEMPLATE].(string), nameTemplate)

		ingress := ingresses.CreateIngress(parsedIngressName, c.namespace, forwardAnnotationsMap, parsedURL,
			parsedURLPath, newServiceObject.Name, getServicePortFromEvent(newServiceObject))

		result, err := c.clientset.ExtensionsV1beta1().Ingresses("first-controller").Create(ingress)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Created ingress", result.GetObjectMeta().GetName())
	} else {
		fmt.Println("Expose label not found, so not creating ingress")
	}
}
func (c *Controller) serviceUpdated(oldObj interface{}, newObj interface{}) {
	fmt.Println("\nService Updated Event")
	c.serviceDeleted(oldObj)
	c.serviceCreated(newObj)
	fmt.Println("Ingress Updated")
}

func (c *Controller) serviceDeleted(deletedServiceObject interface{}) {
	serviceToDelete := deletedServiceObject.(*v1.Service)
	fmt.Println("Service Name: ", serviceToDelete.Name)
	ingressList, err := c.clientset.ExtensionsV1beta1().Ingresses("first-controller").List(meta_v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}

	ingressToRemove := ingresses.GetIngressFromListMatchingGivenServiceName(ingressList, serviceToDelete.Name)
	c.clientset.ExtensionsV1beta1().Ingresses("first-controller").Delete(ingressToRemove.ObjectMeta.Name, &meta_v1.DeleteOptions{})
	fmt.Println("Ingress Deleted with name: ", ingressToRemove.ObjectMeta.Name)
}

func getServicePortFromEvent(service *v1.Service) int {
	return int(service.Spec.Ports[0].Port)
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
		log.Printf("Error syncing ingress %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	log.Printf("Dropping ingress %q out of the queue: %v", key, err)
}
