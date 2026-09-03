package gbcobserveexporterlog

import (
	"fmt"
	"log/slog"
	"strings"

	coreenv "goark.dev/goark/core/env"
)

func newSettings(environment coreenv.Environment, options []Option) (settings, error) {
	resolved := settings{}
	if environment != nil {
		enabled, err := coreenv.ResolveValueAs[bool](environment, "${"+PropertyEnabled+":true}")
		if err != nil {
			return settings{}, err
		}
		spanLevel, err := resolveLevel(environment, PropertySpanLevel, DefaultSpanLevel)
		if err != nil {
			return settings{}, err
		}
		eventLevel, err := resolveLevel(environment, PropertyEventLevel, DefaultEventLevel)
		if err != nil {
			return settings{}, err
		}
		maxAttributes, err := coreenv.ResolveValueAs[int](environment, fmt.Sprintf("${%s:%d}", PropertyMaxAttributes, DefaultMaxAttributes))
		if err != nil {
			return settings{}, err
		}
		resolved.enabled = &enabled
		resolved.spanLevel = &spanLevel
		resolved.eventLevel = &eventLevel
		resolved.maxAttributes = &maxAttributes
	}
	for _, option := range options {
		if option != nil {
			option(&resolved)
		}
	}
	applyDefaults(&resolved)
	if *resolved.maxAttributes <= 0 {
		return settings{}, fmt.Errorf("gbc-observe-exporter-log: max attributes must be positive")
	}
	return resolved, nil
}

func resolveLevel(environment coreenv.Environment, property, fallback string) (slog.Level, error) {
	value, err := coreenv.ResolveValueAs[string](environment, "${"+property+":"+fallback+"}")
	if err != nil {
		return 0, err
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.TrimSpace(value))); err != nil {
		return 0, fmt.Errorf("gbc-observe-exporter-log: invalid %s: %w", property, err)
	}
	return level, nil
}

func applyDefaults(resolved *settings) {
	if resolved.enabled == nil {
		value := DefaultEnabled
		resolved.enabled = &value
	}
	if resolved.spanLevel == nil {
		value := slog.LevelInfo
		resolved.spanLevel = &value
	}
	if resolved.eventLevel == nil {
		value := slog.LevelInfo
		resolved.eventLevel = &value
	}
	if resolved.maxAttributes == nil {
		value := DefaultMaxAttributes
		resolved.maxAttributes = &value
	}
}
