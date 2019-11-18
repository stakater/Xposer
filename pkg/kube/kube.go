package kube

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	routesClient "github.com/openshift/client-go/route/clientset/versioned"
)

func GetRoutesClient() routesClient.Interface {
	config := getClientConfig()
	routesClient, err := routesClient.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Could not create Openshift Routes client: %v", err)
	}

	return routesClient
}

func getClientConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		return getOutOfClusterConfig()
	}

	return config
}

func getOutOfClusterConfig() *rest.Config {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		logrus.Fatalf("Could not get kubernetes config: %v", err)
	}

	return config
}



// GetKubernetesClient returns a k8s clientset to the request from inside of cluster
func GetKubernetesClient() kubernetes.Interface {
	config := getClientConfig()
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Could not create kubernetes client: %v", err)
	}

	return kubeClient
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// IsOpenShift returns true if cluster is openshift based
func IsOpenShift(c kubernetes.Interface) bool {
	client := c.(*kubernetes.Clientset)
	res, err := client.RESTClient().Get().AbsPath("").DoRaw()
	if err != nil {
		return false
	}

	var rp v1.RootPaths
	err = json.Unmarshal(res, &rp)
	if err != nil {
		return false
	}
	for _, p := range rp.Paths {
		if p == "/oapi" {
			return true
		}
	}
	return false
}
