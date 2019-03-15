package config

import (
	"strings"

	"github.com/stakater/Xposer/internal/pkg/constants"
	"k8s.io/api/core/v1"
)

/*
	currentAnnotations contains default config. This method replaces all default config annotations with those annotations provided in
	given service.
*/
func ReplaceDefaultConfigWithProvidedServiceConfig(currentAnnotations map[string]interface{}, serviceObj *v1.Service) map[string]interface{} {
	for annotationKey, annotationValue := range serviceObj.ObjectMeta.Annotations {
		if strings.HasPrefix(annotationKey, constants.INGRESS_CONFIG_ANNOTATION_PREFIX) {
			currentAnnotations[strings.SplitN(annotationKey, "/", 2)[1]] = annotationValue
		}
	}

	return currentAnnotations
}
