package gbcobserveexporterlog

import "log/slog"

type settings struct {
	enabled       *bool
	spanLevel     *slog.Level
	eventLevel    *slog.Level
	maxAttributes *int
}

// Option 定制日志 exporter 自动配置。
type Option func(*settings)

// WithEnabled 显式覆盖日志 exporter 启用状态。
func WithEnabled(enabled bool) Option {
	return func(settings *settings) { settings.enabled = &enabled }
}

// WithSpanLevel 设置 span 摘要日志级别。
func WithSpanLevel(level slog.Level) Option {
	return func(settings *settings) { settings.spanLevel = &level }
}

// WithEventLevel 设置 event 摘要日志级别。
func WithEventLevel(level slog.Level) Option {
	return func(settings *settings) { settings.eventLevel = &level }
}

// WithMaxAttributes 设置单条记录最多携带的观测属性数量。
func WithMaxAttributes(limit int) Option {
	return func(settings *settings) { settings.maxAttributes = &limit }
}
