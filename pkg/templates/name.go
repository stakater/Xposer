package templates

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/stakater/Xposer/pkg/constants"
)

type NameTemplate struct {
	Service   string
	Namespace string
}

func ParseIngressNameTemplate(templateToParse string, nameTemplate *NameTemplate) string {
	var parsedTemplate bytes.Buffer
	fmt.Println("Template to parse: ", templateToParse)

	tmplURL, err := template.New(constants.INGRESS_NAME_TEMPLATE).Parse(templateToParse)
	if err != nil {
		panic(err)
	}
	err = tmplURL.Execute(&parsedTemplate, nameTemplate)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed template: ", parsedTemplate.String())

	return parsedTemplate.String()
}
