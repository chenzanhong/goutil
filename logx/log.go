// logx/setup.go

package logx

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RotateConfig 日志轮转配置
type RotateConfig struct {
	MaxSize    int  // MB
	MaxBackups int  // 备份数
	MaxAge     int  // 天
	Compress   bool // 是否压缩旧日志
}

// DefaultRotateConfig 默认轮转配置
var DefaultRotateConfig = RotateConfig{
	MaxSize:    10,
	MaxBackups: 5,
	MaxAge:     60,
	Compress:   false,
}

// IsZero 判断 RotateConfig 是否为零值
// 如果所有字段都为零值，表示用户不想启用日志轮转
func (r RotateConfig) IsZero() bool {
	return r.MaxSize == 0 && r.MaxBackups == 0 && r.MaxAge == 0 && !r.Compress
}

// SetupZapLogger 创建并返回一个 *zap.Logger
// 如果 rotate 为零值（RotateConfig{}），则不启用日志轮转，直接写入文件
// 否则使用 lumberjack 进行日志轮转
func SetupZapLogger(
	path string,
	level zapcore.Level,
	rotate RotateConfig,
	addCaller, addStacktrace bool,
) (*zap.Logger, error) {
	if err := checkFileWritable(path); err != nil {
		return nil, fmt.Errorf("无法访问日志文件 %s: %v", path, err)
	}

	// 决定日志写入方式
	var writeSyncer zapcore.WriteSyncer
	if rotate.IsZero() {
		// 不启用轮转：直接写入文件
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("无法打开日志文件 %s: %v", path, err)
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		// 启用轮转：使用 lumberjack
		lumberjackLogger := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    rotate.MaxSize,
			MaxBackups: rotate.MaxBackups,
			MaxAge:     rotate.MaxAge,
			Compress:   rotate.Compress,
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
	}

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 构建选项
	var options []zap.Option
	if addCaller {
		options = append(options, zap.AddCaller())
	}
	if addStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger := zap.New(core, options...)
	return logger, nil
}

// 快捷函数：使用默认配置（启用轮转）
func SetupDefaultZapLogger(path string) (*zap.Logger, error) {
	return SetupZapLogger(path, zapcore.InfoLevel, DefaultRotateConfig, true, false)
}

// checkFileWritable 检查路径是否可写（用于早期验证）
func checkFileWritable(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}