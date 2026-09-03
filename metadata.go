package gbcobserveexporterlog

const (
	// ModulePath 是模块的 Go 导入路径。
	ModulePath = "goark.dev/gbc-observe-exporter-log"
	// Repository 是模块对应的官方 Git 仓库名。
	Repository = "goark-boot-contrib-observe-exporter-log"
	// StarterID 是自动配置的稳定标识。
	StarterID = "goark.boot.observe.exporter.log"
	// BeanNameExporter 是日志 exporter Bean 名称。
	BeanNameExporter = "goark.observe.exporter.log"
)

const (
	// PropertyEnabled 控制日志 exporter 是否启用。
	PropertyEnabled = "goark.observe.exporters.log.enabled"
	// PropertySpanLevel 设置 span 摘要日志级别。
	PropertySpanLevel = "goark.observe.exporters.log.span-level"
	// PropertyEventLevel 设置 event 摘要日志级别。
	PropertyEventLevel = "goark.observe.exporters.log.event-level"
	// PropertyMaxAttributes 设置单条日志最多输出的观测属性数量。
	PropertyMaxAttributes = "goark.observe.exporters.log.max-attributes"
)

const (
	// DefaultEnabled 表示显式引入 starter 时默认启用。
	DefaultEnabled = true
	// DefaultSpanLevel 是默认 span 摘要日志级别。
	DefaultSpanLevel = "INFO"
	// DefaultEventLevel 是默认 event 摘要日志级别。
	DefaultEventLevel = "INFO"
	// DefaultMaxAttributes 是单条日志默认允许的最大观测属性数量。
	DefaultMaxAttributes = 64
)
