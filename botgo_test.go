package botgo

import (
	"testing"

	"github.com/tencent-connect/botgo/openapi"
)

func TestUseOpenAPIVersion(t *testing.T) {
	type args struct {
		version openapi.APIVersion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"not found", args{version: 0}, true,
		},
		{
			"v1 found", args{version: 1}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SelectOpenAPIVersion(tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("SelectOpenAPIVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
