package config

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		{
			name: "true",
			args: args{
				filePath: "../../configs/testConfig.yaml",
			},
			want: Configuration{
				Domain: "testDomain.com",
			},
			wantErr: false,
		},
		{
			name: "true",
			args: args{
				filePath: "../../configs/test.yaml",
			},
			want:    Configuration{},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
