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

func ParseIngressURLOrPathTemplate(templateToParse string, URLTemplate *URLTemplate) string {
	var parsedTemplate bytes.Buffer
	logrus.Infof("Template to parse: %v", templateToParse)
	tmplURL, err := template.New("template").Parse(templateToParse)
	if err != nil {
		panic(err)
	}
	err = tmplURL.Execute(&parsedTemplate, URLTemplate)
	if err != nil {
		panic(err)
	}
	logrus.Infof("Parsed template: %v", parsedTemplate.String())
	return parsedTemplate.String()
}
