package gbcobserveexporterlog

import (
	"context"
	"log/slog"

	"goark.dev/boot"
	gbclog "goark.dev/gbc-log"
	gbcobserve "goark.dev/gbc-observe"
	goarkcontainer "goark.dev/goark/container"
	appcontext "goark.dev/goark/context"
	goarklog "goark.dev/log"
	"goark.dev/observe"
	observeexporterlog "goark.dev/observe-exporter-log"
)

const exporterLoggerName = "observe.exporter.log"

// AutoConfigure 创建日志 exporter 自动配置，并补充缺失的日志与观测基础配置。
func AutoConfigure(options ...Option) boot.AutoConfiguration {
	copied := append([]Option(nil), options...)
	return boot.NewAutoConfiguration(StarterID, func(ctx context.Context, app *appcontext.ApplicationContext) error {
		if !hasConfiguration(app, gbclog.StarterID+".configuration") {
			if err := gbclog.AutoConfigure().Configure(ctx, app); err != nil {
				return err
			}
		}
		if !hasConfiguration(app, gbcobserve.StarterID+".configuration") {
			if err := gbcobserve.AutoConfigure().Configure(ctx, app); err != nil {
				return err
			}
		}
		return app.RegisterConfiguration(configuration{options: copied})
	}, boot.WithAutoConfigurationOrder(100))
}

type configuration struct{ options []Option }

func (configuration) Name() string { return StarterID + ".configuration" }
func (configuration) Order() int   { return -15000 }
func (c configuration) Register(ctx context.Context, registry *goarkcontainer.Registry) error {
	return c.RegisterWithContext(ctx, appcontext.NewConfigurationContext(nil, registry))
}
func (c configuration) RegisterWithContext(_ context.Context, config appcontext.ConfigurationContext) error {
	resolved, err := newSettings(config.Environment(), c.options)
	if err != nil {
		return err
	}
	if !*resolved.enabled {
		return nil
	}
	return goarkcontainer.Register[observe.Exporter](config.Registry(), BeanNameExporter, func(ctx context.Context, resolver goarkcontainer.Resolver) (observe.Exporter, error) {
		logger, err := goarkcontainer.Get[*slog.Logger](ctx, resolver, gbclog.BeanNameLogger)
		if err != nil {
			return nil, err
		}
		return observeexporterlog.New(goarklog.WithName(logger, exporterLoggerName),
			observeexporterlog.WithSpanLevel(*resolved.spanLevel),
			observeexporterlog.WithEventLevel(*resolved.eventLevel),
			observeexporterlog.WithMaxAttributes(*resolved.maxAttributes),
		)
	}, goarkcontainer.WithDependsOn(gbclog.BeanNameLogger), goarkcontainer.WithLazy())
}

func hasConfiguration(app *appcontext.ApplicationContext, name string) bool {
	for _, descriptor := range app.Configurations() {
		if descriptor.Name == name {
			return true
		}
	}
	return false
}
