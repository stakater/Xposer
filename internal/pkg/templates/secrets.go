package templates

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"

	"github.com/stakater/Xposer/internal/pkg/constants"
)

type SecretTemplate struct {
	Service   string
	Namespace string
}

func CreateSecretTemplate(service string, namespace string) *SecretTemplate {
	return &SecretTemplate{
		Service:   service,
		Namespace: namespace,
	}
}

func ParseIngressSecretTemplate(templateToParse string, secretTemplate *SecretTemplate) string {
	var parsedTemplate bytes.Buffer
	logrus.Infof("Template to parse: %v", templateToParse)

	tmplURL, err := template.New(constants.SECRET_NAME_TEMPLATE).Parse(templateToParse)
	if err != nil {
		logrus.Errorf("Can not parse the following template : %v, with error: %v", templateToParse, err)
	}
	err = tmplURL.Execute(&parsedTemplate, secretTemplate)
	if err != nil {
		logrus.Errorf("Can not execute template parsing for: %v, with error: %v", secretTemplate, err)
	}
	logrus.Infof("Parsed template: %v", parsedTemplate.String())

	return parsedTemplate.String()
}
