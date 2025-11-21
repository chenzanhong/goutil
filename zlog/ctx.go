package zlog

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
	UserIDKey    ctxKey = "user_id"
	TraceIDKey   ctxKey = "trace_id"
)

func loggerWithContext(ctx context.Context) *zap.Logger {
	logger := Logger()

	var extraFields []zap.Field

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		extraFields = append(extraFields, zap.String("request_id", reqID))
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		extraFields = append(extraFields, zap.String("user_id", userID))
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		extraFields = append(extraFields, zap.String("trace_id", traceID))
	}

	if len(extraFields) > 0 {
		logger = logger.With(extraFields...)
	}
	return logger
}

func sugarWithContext(ctx context.Context) *zap.SugaredLogger {
	logger := Logger()

	var extraFields []zap.Field

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
		extraFields = append(extraFields, zap.String("request_id", reqID))
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		extraFields = append(extraFields, zap.String("user_id", userID))
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		extraFields = append(extraFields, zap.String("trace_id", traceID))
	}

	if len(extraFields) > 0 {
		logger = logger.With(extraFields...)
	}
	return logger.Sugar()
}

func DebugCtx(ctx context.Context, msg string, fields ...Field) {
	loggerWithContext(ctx).Debug(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...Field) {
	loggerWithContext(ctx).Info(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...Field) {
	loggerWithContext(ctx).Warn(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	loggerWithContext(ctx).Error(msg, fields...)
}

func PanicCtx(ctx context.Context, msg string, fields ...Field) {
	loggerWithContext(ctx).Panic(msg, fields...)
}

func FatalCtx(ctx context.Context, msg string, fields ...Field) {
	loggerWithContext(ctx).Fatal(msg, fields...)
}


func DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	sugarWithContext(ctx).Debugf(format, args...)
}

func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	sugarWithContext(ctx).Infof(format, args...)
}

func WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	sugarWithContext(ctx).Warnf(format, args...)
}

func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	sugarWithContext(ctx).Errorf(format, args...)
}

func PanicfCtx(ctx context.Context, format string, args ...interface{}) {
	sugarWithContext(ctx).Panicf(format, args...)
}

func FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	sugarWithContext(ctx).Fatalf(format, args...)
}


func DebugwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	sugarWithContext(ctx).Debugw(msg, keysAndValues...)
}

func InfowCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	sugarWithContext(ctx).Infow(msg, keysAndValues...)
}

func WarnwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	sugarWithContext(ctx).Warnw(msg, keysAndValues...)
}

func ErrorwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	sugarWithContext(ctx).Errorw(msg, keysAndValues...)
}

func PanicwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	sugarWithContext(ctx).Panicw(msg, keysAndValues...)
}

func FatalwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	sugarWithContext(ctx).Fatalw(msg, keysAndValues...)
}
