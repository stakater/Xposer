package services

import "k8s.io/api/core/v1"

func GetServicePortFromEvent(service *v1.Service) int {
	return int(service.Spec.Ports[0].Port)
}
