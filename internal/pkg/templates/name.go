package templates

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"

	"github.com/stakater/Xposer/internal/pkg/constants"
)

type NameTemplate struct {
	Service   string
	Namespace string
}

func CreateNameTemplate(service string, namespace string) *NameTemplate {
	return &NameTemplate{
		Service:   service,
		Namespace: namespace,
	}
}

func ParseIngressNameTemplate(templateToParse string, nameTemplate *NameTemplate) string {
	var parsedTemplate bytes.Buffer
	logrus.Infof("Template to parse: %v", templateToParse)

	tmplURL, err := template.New(constants.INGRESS_NAME_TEMPLATE).Parse(templateToParse)
	if err != nil {
		logrus.Panic(err)
	}
	err = tmplURL.Execute(&parsedTemplate, nameTemplate)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("Parsed template: %v", parsedTemplate.String())

	return parsedTemplate.String()
}
