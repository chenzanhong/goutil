package zlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Global instances (for backward compatibility)
var (
	globalLogger        *zap.Logger
	globalSugaredLogger *zap.SugaredLogger
	once                sync.Once
)

// newLogger creates a new zap.Logger instance with automatic config validation,
// default value filling, and path resolution.
// internal helper, not exported
func newLogger(config LoggerConfig) (*zap.Logger, error) {
	cfg := config

	// Normalize log level
	if cfg.Level < DebugLevel || cfg.Level > FatalLevel {
		cfg.Level = InfoLevel
	}

	// Normalize output destination
	switch cfg.Output {
	case "console", "file", "both":
		// valid
	default:
		cfg.Output = "console"
	}

	// Normalize format
	if cfg.Format != "json" && cfg.Format != "console" {
		cfg.Format = "console"
	}

	// Validate file path when needed
	if (cfg.Output == "file" || cfg.Output == "both") && cfg.FilePath == "" {
		return nil, fmt.Errorf("file path is required when output is 'file' or 'both'")
	}

	// Apply reasonable defaults for rotation settings
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = 100 // MB
	}
	if cfg.MaxBackups < 0 {
		cfg.MaxBackups = 10
	}
	if cfg.MaxAge < 0 {
		cfg.MaxAge = 30 // days
	}

	// Resolve relative file path to absolute
	if cfg.FilePath != "" && !filepath.IsAbs(cfg.FilePath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		cfg.FilePath = filepath.Join(wd, cfg.FilePath)
	}

	// Create log directory if needed
	if cfg.FilePath != "" {
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory %q: %w", dir, err)
		}
	}

	// 4. Build encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 5. Build cores
	var cores []zapcore.Core
	zapLevel := cfg.Level.toZapCoreLevel()

	// Console output
	if cfg.Output == "console" || cfg.Output == "both" {
		var enc zapcore.Encoder
		consoleEncCfg := encoderConfig
		if cfg.Format == "json" {
			enc = zapcore.NewJSONEncoder(consoleEncCfg)
		} else {
			consoleEncCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
			enc = zapcore.NewConsoleEncoder(consoleEncCfg)
		}
		cores = append(cores, zapcore.NewCore(enc, zapcore.Lock(os.Stdout), zapLevel))
	}

	// File output (always JSON)
	if cfg.Output == "file" || cfg.Output == "both" {
		writer := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		var enc zapcore.Encoder
		consoleEncCfg := encoderConfig
		if cfg.Format == "json" {
			enc = zapcore.NewJSONEncoder(consoleEncCfg)
		} else {
			enc = zapcore.NewConsoleEncoder(consoleEncCfg)
		}
		cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(writer), zapLevel))
	}

	if len(cores) == 0 {
		return nil, fmt.Errorf("no valid log output configured")
	}

	// 6. Build logger
	core := zapcore.NewTee(cores...)
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.ErrorOutput(zapcore.Lock(os.Stderr)),
	}

	if cfg.Sampling {
		options = append(options, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
		}))
	}

	logger := zap.New(core, options...)
	if len(cfg.Fields) > 0 {
		var fields []Field
		for k, v := range cfg.Fields {
			fields = append(fields, zap.String(k, v))
		}
		logger = logger.WithOptions(zap.Fields(fields...))
	}

	return logger, nil
}

// InitLogger initializes global logger (thread-safe)
func InitLogger(config LoggerConfig) error {
	var err error
	once.Do(func() {
		globalLogger, err = newLogger(config)
		if err == nil {
			globalSugaredLogger = globalLogger.Sugar()
		}
	})
	return err
}

// Logger returns global zap.Logger
func Logger() *zap.Logger {
	if globalLogger == nil {
		once.Do(func() {
			cfg := DefaultConfig()
			globalLogger, _ = newLogger(cfg)
			globalSugaredLogger = globalLogger.Sugar()
		})
	}
	return globalLogger
}

// Sugar returns global SugaredLogger
func Sugar() *zap.SugaredLogger {
	_ = Logger() // Trigger initialization
	return globalSugaredLogger
}

// InitDefault initializes with default configuration
func InitDefault() error {
	return InitLogger(DefaultConfig())
}

// MustInitDefault panics if default logger fails to initialize.
// Useful in main() for fail-fast behavior.
func MustInitDefault() {
	if err := InitDefault(); err != nil {
		panic(fmt.Sprintf("failed to init default logger: %v", err))
	}
}

// Sync ensures logs are flushed to disk
func Sync() error {
	logger := Logger() // Trigger default initialization if not already initialized
	return logger.Sync()
}
