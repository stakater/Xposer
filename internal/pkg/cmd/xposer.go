package cmd

import (
	"os"

	routeClient "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stakater/Xposer/internal/pkg/config"
	"github.com/stakater/Xposer/internal/pkg/constants"
	"github.com/stakater/Xposer/internal/pkg/controller"
	"github.com/stakater/Xposer/pkg/kube"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewXposerCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "xposer",
		Short: "A Kubernetes controller to watch Services and generate Ingresses/Routes and TLS Certificates automatically",
		Run:   startXposer,
	}
	return cmds
}

func startXposer(cmd *cobra.Command, args []string) {
	currentNamespace := os.Getenv("KUBERNETES_NAMESPACE")
	if currentNamespace == "" {
		currentNamespace = v1.NamespaceAll
		logrus.Infof("KUBERNETES_NAMESPACE is unset, will monitor services in all namespaces.")
	}

	var kubeClient kubernetes.Interface
	var osClient *routeClient.RouteV1Client

	cfg, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = kube.GetClientOutOfCluster()
	} else {
		kubeClient = kube.GetClient()
	}

	var clusterType = constants.KUBERNETES
	if kube.IsOpenShift(kubeClient.(*kubernetes.Clientset)) {
		clusterType = constants.OPENSHIFT
		osClient, err = routeClient.NewForConfig(cfg)
		if err != nil {
			logrus.Errorf("Can not create Openshift client with error: %v", err.Error())
		}
	}

	config := config.GetControllerConfig()
	controller := controller.NewController(kubeClient, osClient, config, clusterType, currentNamespace)

	if currentNamespace != "" {
		logrus.Infof("Controller started in the namespace: %v, with cluster type: %v", currentNamespace, clusterType)
	}

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}
