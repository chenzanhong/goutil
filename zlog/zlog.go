package zlog

import "fmt"

// ========== Structured Logging (High Performance, Recommended for Production) ==========
// Structured logging functions: parameters are []zlog.Field
func Debug(msg string, fields ...Field) {
	executeHooks(DebugLevel, msg, fields)
	Logger().Debug(msg, fields...)
}
func Info(msg string, fields ...Field) {
	executeHooks(InfoLevel, msg, fields)
	Logger().Info(msg, fields...)
}
func Warn(msg string, fields ...Field) {
	executeHooks(WarnLevel, msg, fields)
	Logger().Warn( msg, fields...)
}
func Error(msg string, fields ...Field) {
	executeHooks(ErrorLevel, msg, fields)
	Logger().Error(msg, fields...)
}
func Panic(msg string, fields ...Field) {
	executeHooks(PanicLevel, msg, fields)
	Logger().Panic(msg, fields...)
}
func Fatal(msg string, fields ...Field) {
	executeHooks(FatalLevel, msg, fields)
	Logger().Fatal(msg, fields...)
}

// ========== Key-Value Logging (Easy to Use, Suitable for Rapid Development) ==========
func Debugw(msg string, keysAndValues ...interface{}) {
	executeHooks(DebugLevel, msg, nil)
	Sugar().Debugw(msg, keysAndValues...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	executeHooks(InfoLevel, msg, nil)
	Sugar().Infow(msg, keysAndValues...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	executeHooks(WarnLevel, msg, nil)
	Sugar().Warnw(msg, keysAndValues...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	executeHooks(ErrorLevel, msg, nil)
	Sugar().Errorw(msg, keysAndValues...)
}
func Panicw(msg string, keysAndValues ...interface{}) {
	executeHooks(PanicLevel, msg, nil)
	Sugar().Panicw(msg, keysAndValues...)
}
func Fatalw(msg string, keysAndValues ...interface{}) {
	executeHooks(FatalLevel, msg, nil)
	Sugar().Fatalw(msg, keysAndValues...)
}

// ========== Formatted Logging (fmt Style Compatible) ==========
func Debugf(format string, args ...interface{}) {
	executeHooks(DebugLevel, fmt.Sprintf(format, args...), nil)
	Sugar().Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	executeHooks(InfoLevel, fmt.Sprintf(format, args...), nil)
	Sugar().Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	executeHooks(WarnLevel, fmt.Sprintf(format, args...), nil)
	Sugar().Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	executeHooks(ErrorLevel, fmt.Sprintf(format, args...), nil)
	Sugar().Errorf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	executeHooks(PanicLevel, fmt.Sprintf(format, args...), nil)
	Sugar().Panicf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	executeHooks(FatalLevel, fmt.Sprintf(format, args...), nil)
	Sugar().Fatalf(format, args...)
}
