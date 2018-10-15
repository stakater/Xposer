package ingresses

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stakater/Xposer/internal/pkg/constants"
)

/*
	Looks for Ingress URL Path annotation, and if it doesn't start with a /, appends it
*/
func AppendSlashInPathAnnotationIfNotPresent(currentAnnotations map[string]interface{}) map[string]interface{} {
	if !strings.HasPrefix(currentAnnotations[constants.INGRESS_URL_PATH].(string), "/") {
		currentAnnotations[constants.INGRESS_URL_PATH] = "/" + currentAnnotations[constants.INGRESS_URL_PATH].(string)
	}

	return currentAnnotations
}

/*
	Generate a map of annotations to forward to Ingress
*/
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
