package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/frankgreco/demogods/pkg/metrics"
)

type loggerKey int

const id loggerKey = iota

var logger *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	l, err := config.Build(
		zap.Hooks(addMetrics),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	logger = l
}

// NewContext creates a new context the given contextual fields
func NewContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, id, WithContext(ctx).With(fields...))
}

// WithContext returns a logger from the given context
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}
	if ctxLogger, ok := ctx.Value(id).(*zap.Logger); ok {
		return ctxLogger
	}
	return logger
}

func addMetrics(e zapcore.Entry) error {
	switch e.Level {
	case zap.DebugLevel:
		metrics.LoggingTotal.WithLabelValues("debug").Inc()
	case zap.ErrorLevel:
		metrics.LoggingTotal.WithLabelValues("error").Inc()
	case zap.FatalLevel:
		metrics.LoggingTotal.WithLabelValues("fatal").Inc()
	case zap.InfoLevel:
		metrics.LoggingTotal.WithLabelValues("info").Inc()
	case zap.WarnLevel:
		metrics.LoggingTotal.WithLabelValues("warn").Inc()
	case zap.PanicLevel:
		metrics.LoggingTotal.WithLabelValues("panic").Inc()
	}

	return nil
}
