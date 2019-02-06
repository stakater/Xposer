package configmaps

import (
	"github.com/stakater/Xposer/internal/pkg/constants"

	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateConfigMapObject creates a *v1.Configmap object from given parameters
func CreateConfigMapObject(namespace string, configData map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      constants.XPOSER_CONFIGMAP,
			Namespace: namespace,
		},
		Data: configData,
	}
}

// DeleteFromConfigMapGlobally generates configmap key from given service, and removes that key from xposer configmap from all namespaces
func DeleteFromConfigMapGlobally(clientset kubernetes.Interface, service *v1.Service) {
	namespaces, err := clientset.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	if err != nil {
		logrus.Errorf("Can not fetch all namespaces: %v", err)
	} else {
		for _, namespace := range namespaces.Items {
			configMap, err := clientset.CoreV1().ConfigMaps(namespace.Name).Get(constants.XPOSER_CONFIGMAP, meta_v1.GetOptions{})
			// configmap exist
			if err == nil {
				deleteKeyFromConfigMap(configMap, service, clientset)
			}
		}
	}
}

// DeleteFromConfigMapLocally generates configmap key from given service, and removes that key from xposer configmap in service's namespace
func DeleteFromConfigMapLocally(clientset kubernetes.Interface, service *v1.Service) {
	configMap, err := clientset.CoreV1().ConfigMaps(service.Namespace).Get(constants.XPOSER_CONFIGMAP, meta_v1.GetOptions{})
	// configmap exist
	if err == nil {
		deleteKeyFromConfigMap(configMap, service, clientset)
	}
}

// PopulateConfigMapGlobally creates a new/update existing xposer configmap in all namespaces
func PopulateConfigMapGlobally(clientset kubernetes.Interface, newServiceObject *v1.Service, ingressHost string) {
	namespaces, err := clientset.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	if err != nil {
		logrus.Errorf("Can not fetch all namespaces: %v", err)
	} else {
		for _, namespace := range namespaces.Items {
			configMap, err := clientset.CoreV1().ConfigMaps(namespace.Name).Get(constants.XPOSER_CONFIGMAP, meta_v1.GetOptions{})
			if err != nil {
				createConfigMap(clientset, newServiceObject, ingressHost)
			} else {
				updateConfigMap(configMap, clientset, newServiceObject, ingressHost)
			}
		}
	}
}

// PopulateConfigMapLocally creates a new/update existing xposer configmap in service's namespace
func PopulateConfigMapLocally(clientset kubernetes.Interface, newServiceObject *v1.Service, ingressHost string) {
	configMap, err := clientset.CoreV1().ConfigMaps(newServiceObject.Namespace).Get(constants.XPOSER_CONFIGMAP, meta_v1.GetOptions{})
	if err != nil {
		createConfigMap(clientset, newServiceObject, ingressHost)
	} else {
		updateConfigMap(configMap, clientset, newServiceObject, ingressHost)
	}
}

// createConfigMap uses kubernetes client to create an actual config-map in cluster
func createConfigMap(clientset kubernetes.Interface, newServiceObject *v1.Service, ingressHost string) {
	configData := make(map[string]string)
	configData[newServiceObject.Name+"-"+newServiceObject.Namespace] = ingressHost

	configMap := CreateConfigMapObject(newServiceObject.Namespace, configData)

	_, err := clientset.CoreV1().ConfigMaps(newServiceObject.Namespace).Create(configMap)

	if err != nil {
		logrus.Errorf("Config-map not created in namespace:%v, with error %v", newServiceObject.Namespace, err)
	}

	logrus.Infof("Configmap created in namespace: %v", newServiceObject.Namespace)
}

// updateConfigMap uses kubernetes client to update an actual config-map in cluster
func updateConfigMap(configMap *v1.ConfigMap, clientset kubernetes.Interface, newServiceObject *v1.Service, ingressHost string) {
	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}
	configMap.Data[newServiceObject.Name+"-"+newServiceObject.Namespace] = ingressHost
	_, err := clientset.CoreV1().ConfigMaps(newServiceObject.Namespace).Update(configMap)
	if err != nil {
		logrus.Errorf("Can not update config map in namespace: %v, with error: %v", newServiceObject.Namespace, err)
	}

	logrus.Infof("Configmap updated in namespace: %v", newServiceObject.Namespace)
}

// deleteKeyFromConfigMap uses kubernetes client to delete a key from xposer config-map in cluster
func deleteKeyFromConfigMap(configMap *v1.ConfigMap, service *v1.Service, clientset kubernetes.Interface) {
	delete(configMap.Data, service.Name+"-"+service.Namespace)
	_, err := clientset.CoreV1().ConfigMaps(service.Namespace).Update(configMap)
	if err != nil {
		logrus.Errorf("Can not update config map in namespace: %v, with error: %v", service.Namespace, err)
	}

	logrus.Infof("Configmap updated in namespace: %v", service.Namespace)
}
