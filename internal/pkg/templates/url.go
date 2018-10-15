package templates

import (
	"bytes"
	"html/template"

	"github.com/sirupsen/logrus"
)

type URLTemplate struct {
	Service   string
	Namespace string
	Domain    string
}

func CreateUrlTemplate(service string, namespace string, domain string) *URLTemplate {
	return &URLTemplate{
		Service:   service,
		Namespace: namespace,
		Domain:    domain,
	}
}

func ParseIngressURLOrPathTemplate(templateToParse string, URLTemplate *URLTemplate) string {
	var parsedTemplate bytes.Buffer
	logrus.Infof("Template to parse: %v", templateToParse)
	tmplURL, err := template.New("template").Parse(templateToParse)
	if err != nil {
		logrus.Errorf("Can not parse the following template : %v, with error: %v", templateToParse, err)
	}
	err = tmplURL.Execute(&parsedTemplate, URLTemplate)
	if err != nil {
		logrus.Errorf("Can not execute template parsing for: %v, with error: %v", URLTemplate, err)
	}
	logrus.Infof("Parsed template: %v", parsedTemplate.String())
	return parsedTemplate.String()
}
