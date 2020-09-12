package config

import (
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	os.Setenv("METADATA_FILE", "metadata.json")
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name: "success",
			want: &Config{
				MetadataFile: "metadata.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
