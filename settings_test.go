package gbcobserveexporterlog_test

import (
	"context"
	"testing"

	"goark.dev/boot"
	gbcobserveexporterlog "goark.dev/gbc-observe-exporter-log"
)

func TestAutoConfigureRejectsInvalidAttributeLimit(t *testing.T) {
	app, err := boot.Run(context.Background(), boot.WithAutoConfiguration(
		gbcobserveexporterlog.AutoConfigure(gbcobserveexporterlog.WithMaxAttributes(0)),
	))
	if err == nil {
		_ = app.Close(context.Background())
		t.Fatal("expected invalid attribute limit error")
	}
}
