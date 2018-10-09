package services

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stakater/Xposer/internal/pkg/constants"
	"k8s.io/api/core/v1"
)

/*
	currentAnnotations contains default annoations. This method replaces all default annotations with those annotations provided in
	given service.
*/
func ReplaceAnnotationsInMapWithProvidedServiceAnnotations(currentAnnotations map[string]interface{}, serviceObj *v1.Service) map[string]interface{} {
	for annotationKey, annotationValue := range serviceObj.ObjectMeta.Annotations {
		if strings.HasPrefix(annotationKey, constants.INGRESS_ANNOTATIONS) {
			currentAnnotations[strings.SplitN(annotationKey, "/", 2)[1]] = annotationValue
		}
	}

	return currentAnnotations
}

func AppendsSlashInPathAnnotationIfNotPresent(currentAnnotations map[string]interface{}) map[string]interface{} {
	if !strings.HasPrefix(currentAnnotations[constants.INGRESS_URL_PATH].(string), "/") {
		currentAnnotations[constants.INGRESS_URL_PATH] = "/" + currentAnnotations[constants.INGRESS_URL_PATH].(string)
	}

	return currentAnnotations
}

func FormatURLTemplateAndDeriveURLPath(currentAnnotations map[string]interface{}) map[string]interface{} {
	if strings.Contains(currentAnnotations[constants.INGRESS_URL_TEMPLATE].(string), "/") {
		splittedURLTemplate := strings.SplitN(currentAnnotations[constants.INGRESS_URL_TEMPLATE].(string), "/", 2)
		currentAnnotations[constants.INGRESS_URL_TEMPLATE] = splittedURLTemplate[0]

		if currentAnnotations[constants.INGRESS_URL_PATH].(string) == "/" {
			currentAnnotations[constants.INGRESS_URL_PATH] = currentAnnotations[constants.INGRESS_URL_PATH].(string) + splittedURLTemplate[1]
		}
	}

	return currentAnnotations
}

func CreateForwardAnnotationsMap(splittedAnnotations []string) map[string]string {
	forwardAnnotationsMap := make(map[string]string)

	for _, annotation := range splittedAnnotations {
		parsedAnnotation := strings.Split(annotation, ":")
		if len(parsedAnnotation) != 2 {
			logrus.Warningf("Wrong annotation provided to forward to ingress : %v", annotation)
		} else {
			forwardAnnotationsMap[parsedAnnotation[0]] = strings.Trim(parsedAnnotation[1], " ")
		}
	}

	return forwardAnnotationsMap
}
