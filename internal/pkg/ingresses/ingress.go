package ingresses

import (
	"fmt"

	"github.com/stakater/Xposer/internal/pkg/constants"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetIngressFromListMatchingGivenServiceName(ingressList *v1beta1.IngressList, serviceName string) v1beta1.Ingress {
	var matchedIngress v1beta1.Ingress

	for _, ingress := range ingressList.Items {
		if ingress.Spec.Backend.ServiceName == serviceName {
			fmt.Println("Ingress & Service Name matched")
			matchedIngress = ingress
			break
		}
	}

	return matchedIngress
}

func AddTLSInfoToIngress(ingress v1beta1.Ingress, ingressName string, ingressHost string) v1beta1.Ingress {
	ingress.Spec.TLS = []v1beta1.IngressTLS{
		v1beta1.IngressTLS{
			Hosts:      []string{ingressHost},
			SecretName: ingressName + constants.CERT,
		},
	}

	return ingress
}

func CreateIngress(parsedIngressName string, namespace string, forwardAnnotationsMap map[string]string,
	ingressHost string, ingressPath string, serviceName string, servicePort int) *v1beta1.Ingress {

	return &v1beta1.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        parsedIngressName,
			Namespace:   namespace,
			Annotations: forwardAnnotationsMap,
		},
		Spec: v1beta1.IngressSpec{
			Backend: &v1beta1.IngressBackend{
				ServiceName: serviceName,
				ServicePort: intstr.FromInt(servicePort),
			},
			Rules: []v1beta1.IngressRule{
				v1beta1.IngressRule{
					Host: ingressHost,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								v1beta1.HTTPIngressPath{
									Path: ingressPath,
									Backend: v1beta1.IngressBackend{
										ServiceName: serviceName,
										ServicePort: intstr.FromInt(servicePort),
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
