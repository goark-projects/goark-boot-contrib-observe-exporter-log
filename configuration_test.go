package gbcobserveexporterlog_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"goark.dev/boot"
	gbclog "goark.dev/gbc-log"
	gbcobserve "goark.dev/gbc-observe"
	gbcobserveexporterlog "goark.dev/gbc-observe-exporter-log"
	"goark.dev/goark"
	coreenv "goark.dev/goark/core/env"
	"goark.dev/log"
	"goark.dev/observe"
)

func TestAutoConfigureExportsSpansAndEventsThroughGoarkLog(t *testing.T) {
	var output bytes.Buffer
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(
		gbclog.AutoConfigure(gbclog.WithLoggerContextFactory(testLoggerFactory(&output))),
		gbcobserveexporterlog.AutoConfigure(),
	))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, _ := app.Context()
	provider := goark.MustGet[observe.Provider](t.Context(), appContext, gbcobserve.BeanNameProvider)
	ctx, span := provider.Tracer("starter.test").Start(t.Context(), "admin.request")
	provider.Eventer("starter.test").Emit(ctx, "admin.event", observe.String("route", "/admin/users/{id}"))
	span.End()
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	logs := output.String()
	for _, expected := range []string{"observability span completed", "observability event emitted", "admin.request", "/admin/users/{id}", "trace_id", "goark.dev.observe.exporter.log"} {
		if !strings.Contains(logs, expected) {
			t.Fatalf("logs do not contain %q: %s", expected, logs)
		}
	}
}

func TestAutoConfigureDisabledDoesNotExport(t *testing.T) {
	var output bytes.Buffer
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(
		gbclog.AutoConfigure(gbclog.WithLoggerContextFactory(testLoggerFactory(&output))),
		gbcobserveexporterlog.AutoConfigure(gbcobserveexporterlog.WithEnabled(false)),
	))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, _ := app.Context()
	provider := goark.MustGet[observe.Provider](t.Context(), appContext, gbcobserve.BeanNameProvider)
	_, span := provider.Tracer("starter.test").Start(t.Context(), "disabled.span")
	span.End()
	if _, err := appContext.Get(t.Context(), gbcobserveexporterlog.BeanNameExporter); err == nil {
		t.Fatal("disabled configuration registered exporter bean")
	}
	logger := goark.MustGet[*slog.Logger](t.Context(), appContext, gbclog.BeanNameLogger)
	logger.InfoContext(t.Context(), "ordinary application log")
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	logs := output.String()
	if !strings.Contains(logs, "ordinary application log") {
		t.Fatalf("ordinary log is missing: %s", logs)
	}
	if strings.Contains(logs, "disabled.span") {
		t.Fatalf("disabled exporter emitted telemetry: %s", logs)
	}
}

func TestProviderShutdownRunsBeforeLoggerClose(t *testing.T) {
	var output bytes.Buffer
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(
		gbclog.AutoConfigure(gbclog.WithLoggerContextFactory(testLoggerFactory(&output))),
		gbcobserve.AutoConfigure(gbcobserve.WithExporters(shutdownLoggingExporter{})),
		gbcobserveexporterlog.AutoConfigure(),
	))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if logs := output.String(); !strings.Contains(logs, "exporter shutdown log") {
		t.Fatalf("provider shutdown could not write through goark-log: %s", logs)
	}
}

type shutdownLoggingExporter struct{}

func (shutdownLoggingExporter) Descriptor() observe.ExporterDescriptor {
	return observe.ExporterDescriptor{
		Name:         "test.shutdown-logging",
		Signals:      observe.SignalTraces,
		Stability:    observe.StabilityStable,
		Capabilities: observe.ExporterCapabilities{Push: true},
	}
}
func (shutdownLoggingExporter) ForceFlush(context.Context) error { return nil }
func (shutdownLoggingExporter) Shutdown(ctx context.Context) error {
	slog.Default().InfoContext(ctx, "exporter shutdown log")
	return nil
}
func (shutdownLoggingExporter) ExportSpans(context.Context, []observe.SpanSnapshot) error {
	return nil
}

func testLoggerFactory(output *bytes.Buffer) gbclog.LoggerContextFactory {
	return func(context.Context, coreenv.Environment) (*log.LoggerContext, error) {
		layout, err := log.NewPatternLayout("%logger : %msg%attrs%n")
		if err != nil {
			return nil, err
		}
		return log.NewLoggerContext(log.Options{
			Appenders: []log.Appender{log.NewConsoleAppender(
				log.WithConsoleWriter(output),
				log.WithConsoleLayout(layout),
			)},
			Root: log.RootLogger{Level: slog.LevelInfo, AppenderRefs: []string{"console"}},
		})
	}
}
