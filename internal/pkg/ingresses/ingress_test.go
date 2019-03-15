package ingresses

import (
	"testing"

	"github.com/stakater/Xposer/internal/pkg/constants"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsEmpty(t *testing.T) {
	type args struct {
		ingress v1beta1.Ingress
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Is empty test, should return true",
			args: args{},
			want: true,
		},
		{
			name: "Is empty test, should return false",
			args: args{
				ingress: createIngressWithName(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.ingress); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTLSInfo(t *testing.T) {
	type args struct {
		ingress     *v1beta1.Ingress
		ingressName string
		ingressHost string
	}
	tests := []struct {
		name     string
		args     args
		modified *v1beta1.Ingress
	}{
		{
			name: "Should Add TLS",
			args: args{
				ingress:     createIngressForTLS(),
				ingressHost: "test-host",
				ingressName: "test-ingress-name",
			},
			modified: createIngressForTLS(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddTLSInfo(tt.args.ingress, tt.args.ingressName, tt.args.ingressHost)
			if len(tt.args.ingress.Spec.TLS) < 1 {
				t.Errorf("TLS Not added to ingress = %v", tt.args.ingress)
			}
		})
	}
}

func TestShouldAddTLS(t *testing.T) {
	type args struct {
		ingressConfig map[string]interface{}
		defaultTLS    bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should retrun true",
			args: args{
				ingressConfig: createIngressTlsConfigMap("true"),
			},
			want: true,
		},
		{
			name: "should retrun true",
			args: args{
				ingressConfig: createIngressTlsConfigMap("false"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldAddTLS(tt.args.ingressConfig, tt.args.defaultTLS); got != tt.want {
				t.Errorf("ShouldAddTLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddDefaultBackend(t *testing.T) {
	type args struct {
		ingress     *v1beta1.Ingress
		serviceName string
		servicePort int
	}
	tests := []struct {
		name     string
		args     args
		modified *v1beta1.Ingress
		want     string
	}{
		{
			name: "Should Add Default Backend",
			args: args{
				ingress:     createIngressForDefaultBackend(),
				serviceName: "test-service",
				servicePort: 1234,
			},
			modified: createIngressForDefaultBackend(),
			want: "test-service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddDefaultBackend(tt.args.ingress, tt.args.serviceName, tt.args.servicePort)
			if got := (tt.args.ingress.Spec.Backend.ServiceName); got != tt.want {
				t.Errorf("Default Backend Not added to ingress = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldAddDefaultBackend(t *testing.T) {
	type args struct {
		ingressConfig            map[string]interface{}
		defaultDefaultBackend    bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true",
			args: args{
				ingressConfig: createIngressDefaultBackendConfigMap("true"),
			},
			want: true,
		},
		{
			name: "should return false",
			args: args{
				ingressConfig: createIngressDefaultBackendConfigMap("false"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldAddDefaultBackend(tt.args.ingressConfig, tt.args.defaultDefaultBackend); got != tt.want {
				t.Errorf("ShouldAddDefaultBackend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createIngressWithName() v1beta1.Ingress {
	return v1beta1.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "Test-Ingress",
		},
	}
}

func createIngressForTLS() *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "TLS-Ingress",
		},
	}
}

func createIngressForDefaultBackend() *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "DefaultBackend-Ingress",
		},
	}
}

func createIngressTlsConfigMap(tlsValue string) map[string]interface{} {
	ingressConfig := make(map[string]interface{})
	ingressConfig[constants.TLS] = tlsValue

	return ingressConfig
}

func createIngressDefaultBackendConfigMap(defaultBackendValue string) map[string]interface{} {
	ingressConfig := make(map[string]interface{})
	ingressConfig[constants.DEFAULTBACKEND] = defaultBackendValue

	return ingressConfig
}
