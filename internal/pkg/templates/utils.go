package templates

import (
	"strings"

	"github.com/stakater/Xposer/internal/pkg/constants"
)

/*
	Looks for / in URL-Template; if present it divides the template in 2 parts, and uses the 1st part
	as URL
*/
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
