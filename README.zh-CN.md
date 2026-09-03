# Goark Boot 观测日志导出器

[English](README.md)

`goark.dev/gbc-observe-exporter-log` 为 Goark Boot 自动配置 `goark.dev/observe-exporter-log`，把 span 和观测事件输出到 `goark.dev/gbc-log` 提供的主 `*slog.Logger`。

## 职责边界

- 依赖方向固定为 `observe -> observe-exporter-log -> slog -> goark-log`。
- `goark.dev/log` 永远不依赖观测体系。
- 只导出 span 和 event；不导出 log，避免递归；不导出 metric，避免日志洪泛。
- exporter 不持有 Logger 生命周期。Boot 在普通组件停止后、Provider 关闭前排空 goark-log。

## 使用

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(
    gbcobserveexporterlog.AutoConfigure(),
))
```

缺少基础 starter 时，本 starter 会自动补充 `gbc-log` 和 `gbc-observe`。需要自定义基础 starter 参数时，应先显式注册对应基础 starter。

## Bean

- `goark.observe.exporter.log`：延迟创建的 `observe.Exporter`，由 `gbc-observe` 创建 Provider 时统一收集。

## 配置

- `goark.observe.exporters.log.enabled`：默认为 `true`。
- `goark.observe.exporters.log.span-level`：默认为 `INFO`。
- `goark.observe.exporters.log.event-level`：默认为 `INFO`。
- `goark.observe.exporters.log.max-attributes`：默认为 `64`，且必须为正数。

## 生产注意事项

- 配置采样策略，限制 span 日志量。
- 禁止把这些日志再次送入观测日志桥接器。
- 禁止把密钥、原始 URL、SQL 参数、用户标识等敏感或高基数数据放入观测属性。
- 外部传入的 Provider 自行负责 exporter 装配，不收集 Boot exporter Bean。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
```

本项目采用 Apache License 2.0。
