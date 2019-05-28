package ingresses

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/stakater/Xposer/internal/pkg/config"
	"github.com/stakater/Xposer/internal/pkg/constants"
	"github.com/stakater/Xposer/internal/pkg/services"
	"github.com/stakater/Xposer/internal/pkg/templates"
)

type IngressInfo struct {
	IngressName           string
	Namespace             string
	ForwardAnnotationsMap map[string]string
	IngressHost           string
	IngressPath           string
	ServiceName           string
	ServicePort           int
	AddTLS                bool
	SecretName            string
}

func CreateIngressInfo(newServiceObject *v1.Service, configuration config.Configuration) IngressInfo {
	splittedAnnotations := strings.Split(string(newServiceObject.ObjectMeta.Annotations[constants.FORWARD_ANNOTATION]), "\n")
	forwardAnnotationsMap := make(map[string]string)
	ingressConfig := structs.Map(configuration)

	// Overrides default annotains with annotations from new service object
	ingressConfig = config.ReplaceDefaultConfigWithProvidedServiceConfig(ingressConfig, newServiceObject)

	// Adds "/" in URL Path, if user has entered path annotaion without "/"
	ingressConfig = AppendSlashInPathAnnotationIfNotPresent(ingressConfig)

	//	Removes the content after "/" from URL-Template, and if user has not specified path from annotation, use the content after "/" as URL-Path
	ingressConfig = templates.FormatURLTemplateAndDeriveURLPath(ingressConfig)

	// Creates a map of annotations to forward to Ingress
	forwardAnnotationsMap = CreateForwardAnnotationsMap(splittedAnnotations)

	// Generates URL Templates to parse Xposer Specific Annotations
	urlTemplate := templates.CreateUrlTemplate(newServiceObject.Name, newServiceObject.Namespace, ingressConfig[constants.DOMAIN].(string))
	nameTemplate := templates.CreateNameTemplate(newServiceObject.Name, newServiceObject.Namespace)

	// Generate Secret Template to create Secrets
	secretTemplate := templates.CreateSecretTemplate(newServiceObject.Name, newServiceObject.Namespace)

	parsedURL := templates.ParseIngressURLOrPathTemplate(ingressConfig[constants.INGRESS_URL_TEMPLATE].(string), urlTemplate)
	parsedURLPath := templates.ParseIngressURLOrPathTemplate(ingressConfig[constants.INGRESS_URL_PATH].(string), urlTemplate)
	parsedIngressName := templates.ParseIngressNameTemplate(ingressConfig[constants.INGRESS_NAME_TEMPLATE].(string), nameTemplate)
	parsedSecret := templates.ParseIngressSecretTemplate(ingressConfig[constants.SECRET_NAME_TEMPLATE].(string), secretTemplate)

	return IngressInfo{
		IngressName:           parsedIngressName,
		Namespace:             newServiceObject.Namespace,
		ForwardAnnotationsMap: forwardAnnotationsMap,
		IngressHost:           parsedURL,
		IngressPath:           parsedURLPath,
		ServiceName:           newServiceObject.Name,
		ServicePort:           services.GetServicePortFromEvent(newServiceObject),
		AddTLS:                ShouldAddTLS(ingressConfig, configuration.TLS),
		SecretName:            parsedSecret,
	}
}
