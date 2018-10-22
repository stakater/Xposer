package ingresses

import (
	"testing"
)

func TestCreateForwardAnnotationsMap(t *testing.T) {
	type args struct {
		splittedAnnotations []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should add annotations with multiple colons",
			args: args{
				splittedAnnotations: []string{"hello:hello:hello:hello"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateForwardAnnotationsMap(tt.args.splittedAnnotations)
			if len(got) < 1 {
				t.Errorf("Annotations not copied properly")
			}
		})
	}
}
