# Goark Boot Observe Log Exporter

[简体中文](README.zh-CN.md)

`goark.dev/gbc-observe-exporter-log` auto-configures `goark.dev/observe-exporter-log` for Goark Boot. It connects observability spans and events to the primary `*slog.Logger` supplied by `goark.dev/gbc-log`.

## Boundary

- Dependency direction is `observe -> observe-exporter-log -> slog -> goark-log`.
- `goark.dev/log` never depends on observability.
- Only spans and events are exported. Logs are excluded to prevent recursion, and metrics are excluded to prevent log floods.
- The exporter does not own the logger. Boot lifecycle drains goark-log after ordinary components and before the observability provider shuts down.

## Usage

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(
    gbcobserveexporterlog.AutoConfigure(),
))
```

The starter adds `gbc-log` and `gbc-observe` automatically when they are absent. Register either base starter explicitly first when custom options are required.

## Bean

- `goark.observe.exporter.log`: lazy `observe.Exporter`, collected by `gbc-observe` during provider creation.

## Properties

- `goark.observe.exporters.log.enabled`: defaults to `true`.
- `goark.observe.exporters.log.span-level`: defaults to `INFO`.
- `goark.observe.exporters.log.event-level`: defaults to `INFO`.
- `goark.observe.exporters.log.max-attributes`: defaults to `64` and must be positive.

## Production Notes

- Configure sampling to bound span-log volume.
- Never route these records back into an observability log bridge.
- Do not attach secrets, raw URLs, SQL arguments, user identifiers, or other sensitive/high-cardinality values to telemetry.
- An externally supplied Provider owns its exporter wiring and does not collect Boot exporter beans.

## Validation

```bash
go test ./...
go test -race ./...
go vet ./...
```

Licensed under Apache License 2.0.
