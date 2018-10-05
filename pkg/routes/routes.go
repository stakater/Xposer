package routes

import (
	osV1 "github.com/openshift/api/route/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateRoute(routeName string, namespace string, forwardAnnotationsMap map[string]string,
	routeHost string, routePath string, serviceName string, servicePort int) *osV1.Route {
	return &osV1.Route{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        routeName,
			Namespace:   namespace,
			Annotations: forwardAnnotationsMap,
		},
		Spec: osV1.RouteSpec{
			Host: routeHost,
			Path: routePath,
			To: osV1.RouteTargetReference{
				Kind: "Service",
				Name: serviceName,
			},
			Port: &osV1.RoutePort{
				TargetPort: intstr.FromInt(servicePort),
			},
		},
	}
}
