package gbcobserveexporterlog_test

import (
	"testing"

	gbcobserveexporterlog "goark.dev/gbc-observe-exporter-log"
)

func TestModuleMetadata(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "module path", got: gbcobserveexporterlog.ModulePath, want: "goark.dev/gbc-observe-exporter-log"},
		{name: "repository", got: gbcobserveexporterlog.Repository, want: "goark-boot-contrib-observe-exporter-log"},
		{name: "starter ID", got: gbcobserveexporterlog.StarterID, want: "goark.boot.observe.exporter.log"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("value = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
