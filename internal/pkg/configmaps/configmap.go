package configmaps

import (
	"github.com/stakater/Xposer/internal/pkg/constants"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(namespace string, configData map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      constants.XPOSER_CONFIGMAP,
			Namespace: namespace,
		},
		Data: configData,
	}
}
