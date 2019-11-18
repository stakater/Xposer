package cmd

import (
	"os"

	routesClient "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stakater/Xposer/internal/pkg/config"
	"github.com/stakater/Xposer/internal/pkg/constants"
	"github.com/stakater/Xposer/internal/pkg/controller"
	"github.com/stakater/Xposer/pkg/kube"
	"k8s.io/api/core/v1"
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

	var routesClient routesClient.Interface

	kubeClient := kube.GetKubernetesClient()

	var clusterType = constants.KUBERNETES
	if kube.IsOpenShift(kubeClient) {
		clusterType = constants.OPENSHIFT
		routesClient = kube.GetRoutesClient()
	}

	config := config.GetControllerConfig()
	controller := controller.NewController(kubeClient, routesClient, config, clusterType, currentNamespace)

	if currentNamespace != "" {
		logrus.Infof("Controller started in the namespace: %v, with cluster type: %v", currentNamespace, clusterType)
	}

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}
