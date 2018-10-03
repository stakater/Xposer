package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

type URLTemplate struct {
	Service   string
	Namespace string
	Domain    string
}

func ParseIngressURLOrPathTemplate(templateToParse string, URLTemplate *URLTemplate) string {
	var parsedTemplate bytes.Buffer
	fmt.Println("Template to parse: ", templateToParse)
	tmplURL, err := template.New("template").Parse(templateToParse)
	if err != nil {
		panic(err)
	}
	err = tmplURL.Execute(&parsedTemplate, URLTemplate)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed template: ", parsedTemplate.String())
	return parsedTemplate.String()
}
