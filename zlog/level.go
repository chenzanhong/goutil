package zlog

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Level represents log level, hiding zapcore.Level internally
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	PanicLevel Level = "panic"
	FatalLevel Level = "fatal"
)

// String returns human-readable level name
func (l Level) String() string {
	return string(l)
}

// Valid checks if the level is one of the predefined valid levels.
func (l Level) Valid() bool {
	switch l {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel, PanicLevel, FatalLevel:
		return true
	default:
		return false
	}
}

// UnmarshalText implements encoding.TextUnmarshaler
// Supports parsing from YAML, JSON, TOML, env vars, etc.
func (l *Level) UnmarshalText(text []byte) error {
	levelStr := strings.ToLower(string(text))
	switch levelStr {
	case "debug", "d":
		*l = DebugLevel
	case "info", "i":
		*l = InfoLevel
	case "warn", "warning", "w":
		*l = WarnLevel
	case "error", "err", "e":
		*l = ErrorLevel
	case "panic", "p":
		*l = PanicLevel
	case "fatal", "f":
		*l = FatalLevel
	default:
		return fmt.Errorf("invalid log level: %q", string(text))
	}
	return nil
}

func (l Level) MarshalText() ([]byte, error) {
	if !l.Valid() {
		return []byte("info"), nil // safe default
	}
	return []byte(l), nil
}

// toZapCoreLevel converts to zapcore.Level (internal use)
func (l Level) toZapCoreLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// fromZapCoreLevel converts from zapcore.Level (if needed)
func fromZapCoreLevel(l zapcore.Level) Level {
	switch l {
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.WarnLevel:
		return WarnLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.PanicLevel:
		return PanicLevel
	case zapcore.FatalLevel:
		return FatalLevel
	default:
		return InfoLevel
	}
}
