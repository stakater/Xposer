package ingresses

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/stakater/Xposer/internal/pkg/constants"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateFromIngressInfo(ingresInfo IngressInfo) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        ingresInfo.IngressName,
			Namespace:   ingresInfo.Namespace,
			Annotations: ingresInfo.ForwardAnnotationsMap,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				v1beta1.IngressRule{
					Host: ingresInfo.IngressHost,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								v1beta1.HTTPIngressPath{
									Path: ingresInfo.IngressPath,
									Backend: v1beta1.IngressBackend{
										ServiceName: ingresInfo.ServiceName,
										ServicePort: intstr.FromInt(ingresInfo.ServicePort),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func IsEmpty(ingress v1beta1.Ingress) bool {
	if ingress.Name == "" {
		return true
	}

	return false
}

func GetFromListMatchingGivenServiceName(ingressList *v1beta1.IngressList, serviceName string) v1beta1.Ingress {
	var matchedIngress v1beta1.Ingress

	for _, ingress := range ingressList.Items {

		if ingress.Spec.Backend != nil && ingress.Spec.Backend.ServiceName == serviceName {
			matchedIngress = ingress
			break
		}
	}

	return matchedIngress
}

func AddDefaultBackend(ingress *v1beta1.Ingress, serviceName string, servicePort int) {
	ingress.Spec.Backend = &v1beta1.IngressBackend{
		ServiceName: serviceName,
		ServicePort: intstr.FromInt(servicePort),
	}
}

func AddTLSInfo(ingress *v1beta1.Ingress, ingressName string, ingressHost string) {
	ingress.Spec.TLS = []v1beta1.IngressTLS{
		v1beta1.IngressTLS{
			Hosts:      []string{ingressHost},
			SecretName: ingressName + constants.CERT,
		},
	}
}

func ShouldAddDefaultBackend(ingressConfig map[string]interface{}, defaultDefaultBackend bool) bool {
	switch defaultBackendSwitch := ingressConfig[constants.DEFAULTBACKEND].(type) {
	case string:
		defaultBackend, err := strconv.ParseBool(defaultBackendSwitch)
		if err != nil {
			logrus.Warnf("The value of Default Backend annotation is wrong. It should only be true or false. Reverting to default value: %v", defaultDefaultBackend)
			ingressConfig[constants.DEFAULTBACKEND] = defaultDefaultBackend
		} else {
			ingressConfig[constants.DEFAULTBACKEND] = defaultBackend
		}
		break
	}

	if ingressConfig[constants.DEFAULTBACKEND] == true {
		return true
	}

	return false
}

func ShouldAddTLS(ingressConfig map[string]interface{}, defaultTLS bool) bool {
	switch tlsSwitch := ingressConfig[constants.TLS].(type) {
	case string:
		tls, err := strconv.ParseBool(tlsSwitch)
		if err != nil {
			logrus.Warnf("The value of TLS annotation is wrong. It should only be true or false. Reverting to default value: %v", defaultTLS)
			ingressConfig[constants.TLS] = defaultTLS
		} else {
			ingressConfig[constants.TLS] = tls
		}
		break
	}

	if ingressConfig[constants.TLS] == true {
		return true
	}

	return false
}
