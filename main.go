package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/stakater/Xposer/pkg/config"
	"github.com/stakater/Xposer/pkg/controller"
	"github.com/stakater/Xposer/pkg/kube"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	currentNamespace := os.Getenv("KUBERNETES_NAMESPACE")
	if len(currentNamespace) == 0 {
		currentNamespace = v1.NamespaceAll
		log.Println("Warning: KUBERNETES_NAMESPACE is unset, will monitor ingresses in all namespaces.")
	}

	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = kube.GetClientOutOfCluster()
	} else {
		kubeClient = kube.GetClient()
	}

	var resource = "services"
	if kube.IsOpenShift(kubeClient.(*kubernetes.Clientset)) {
		resource = "routes"
	}

	config := getControllerConfig()
	fmt.Println("Config: ", config)
	// Now let's start the controller
	fmt.Println("Initializing Controller")
	controller := controller.NewController(kubeClient, config, resource, currentNamespace)
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}

func getClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	//If kube config file exists in home so use that
	if _, err := os.Stat(kubeconfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		//use Incluster Configuration
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
func getControllerConfig() config.Configuration {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "configs/config.yaml"
	}

	configuration := config.ReadConfig(configFilePath)
	return configuration
}
