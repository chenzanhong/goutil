# zlog

zlog是一个基于zap的高性能Go语言日志库，提供简洁易用的API和丰富的配置选项，适用于各种Go应用程序的日志记录需求。

## 特性

- 高性能：基于uber-go/zap构建，保持卓越性能
- 多种日志风格：支持结构化日志、键值对日志和格式化日志
- 灵活配置：支持控制台输出、文件输出或同时输出
- 自动轮转：支持日志文件自动轮转和压缩
- 环境变量配置：支持通过环境变量进行配置
- 日志钩子：支持自定义日志钩子进行扩展
- 线程安全：全局实例的初始化是线程安全的

## 安装

```bash
# 在Go项目中添加本地依赖
go get github.com/chenzanhong/zlog
```

## 快速开始

### 默认初始化

最简单的方式是使用默认配置初始化日志系统：

```go
import (
    "github.com/chenzanhong/zlog"
)

func main() {
    // Initialize logging system with default configuration
    err := zlog.InitLoggerDefault()
    if err != nil {
        panic("初始化日志失败: " + err.Error())
    }
    defer zlog.Sync() // Ensure logs are flushed to disk
    
    // Use logging
    zlog.Info("应用启动", zlog.String("version", "1.0.0"))
    zlog.Infow("用户登录", "username", "admin", "ip", "127.0.0.1")
    zlog.Infof("处理请求耗时: %v", 100*time.Millisecond)
}
```

### 自定义配置初始化

```go
import (
    "github.com/chenzanhong/zlog"
)

func main() {
    config := &zlog.LoggerConfig{
        Level:      "debug",          // Log level: debug, info, warn, error, panic, fatal
        Output:     "both",           // Output target: console, file, both
        Format:     "json",           // Console format: json, console; file is always json
        FilePath:   "./logs/app.log", // Log file path
        MaxSize:    100,              // Max size per log file (MB)
        MaxBackups: 10,               // Max number of backup log files
        MaxAge:     30,               // Max days to retain log files
        Compress:   true,             // Whether to compress old log files
        Sampling:   false,            // Whether to enable log sampling
    }
    
    err := zlog.InitLogger(config)
    if err != nil {
        panic("初始化日志失败: " + err.Error())
    }
    defer zlog.Sync()
    
    // Use logging
    zlog.Debug("调试信息", zlog.Int("count", 10))
}
```

## 配置说明

### LoggerConfig 结构体

| 字段名      | 类型   | 默认值      | 说明                              | 环境变量         |
|----------|------|----------|---------------------------------|--------------|
| Level    | string | "info"  | 日志级别                            | LOG_LEVEL    |
| Output   | string | "both"  | 输出目标：console, file, both       | LOG_OUTPUT   |
| Format   | string | "console" | 控制台格式：json, console            | LOG_FORMAT   |
| FilePath | string | "./logs/app.log" | 日志文件路径                          | LOG_FILE_PATH |
| MaxSize  | int  | 100      | 单个日志文件最大大小(MB)                  | LOG_MAX_SIZE |
| MaxBackups | int  | 10       | 保留的最大日志文件数                      | LOG_MAX_BACKUPS |
| MaxAge   | int  | 30       | 保留的最大天数                         | LOG_MAX_AGE  |
| Compress | bool | true     | 是否压缩旧日志文件                       | LOG_COMPRESS |
| Sampling | bool | false    | 是否启用日志采样                        | LOG_SAMPLING |

## 使用指南

### 结构化日志（推荐生产环境使用）

结构化日志使用 `Field` 类型来记录键值对，性能更高，适合生产环境：

```go
// Basic usage
zlog.Debug("调试信息")
zlog.Info("普通信息")
zlog.Warn("警告信息")
zlog.Error("错误信息")
zlog.Panic("panic message") // Will trigger panic
zlog.Fatal("fatal error") // Will cause program exit

// Usage with fields
zlog.Info("用户登录成功", 
    zlog.String("username", "admin"),
    zlog.String("ip", "127.0.0.1"),
    zlog.Int("user_id", 1001),
    zlog.Bool("success", true),
    zlog.Time("login_time", time.Now()),
)
```

### 键值对日志（易用，适合快速开发）

键值对日志使用更简单的键值对形式：

```go
zlog.Debugw("调试信息", "key1", "value1", "key2", 2)
zlog.Infow("用户操作", "user", "admin", "action", "create", "id", 100)
zlog.Warnw("资源警告", "resource", "database", "usage", "90%")
zlog.Errorw("API调用失败", "endpoint", "/api/users", "status", 500, "latency", "100ms")
```

### 格式化日志（兼容 fmt 风格）

格式化日志使用类似 fmt.Printf 的风格：

```go
zlog.Debugf("当前时间: %v", time.Now())
zlog.Infof("处理请求 %s，耗时 %v", requestID, duration)
zlog.Warnf("资源使用率过高: %.2f%%", usageRate)
zlog.Errorf("连接数据库失败: %v", err)
```

### 日志字段类型

zlog提供了多种字段类型用于结构化日志：

| 函数名      | 类型      | 示例                          |
|----------|---------|-----------------------------|
| String   | string  | `zlog.String("name", "admin")` |
| Int      | int     | `zlog.Int("count", 10)`        |
| Int64    | int64   | `zlog.Int64("size", 1024)`     |
| Bool     | bool    | `zlog.Bool("success", true)`   |
| Float64  | float64 | `zlog.Float64("rate", 0.95)`   |
| Duration | time.Duration | `zlog.Duration("latency", time.Millisecond*100)` |
| Time     | time.Time | `zlog.Time("timestamp", time.Now())` |
| Any      | interface{} | `zlog.Any("data", user)`     |

### 日志钩子

zlog支持自定义日志钩子，可以在日志记录时执行额外的操作：

```go
import (
    "github.com/chenzanhong/zlog"
)

// Custom log hooks
type AlertHook struct{}

func (h *AlertHook) OnLog(level zlog.Level, msg string, fields []zlog.Field) error {
    // Example: Send alert when error level logs appear
    if level >= zlog.ErrorLevel {
        // Alert sending logic
        // ...
    }
    return nil
}

func main() {
    // Initialize logging
    err := zlog.InitLoggerDefault()
    if err != nil {
        panic(err)
    }
    
    // Register log hook
    zlog.RegisterLogHook(&AlertHook{})
    
    // Use logging
    zlog.Error("这是一个错误", zlog.String("reason", "测试"))
}
```

## 最佳实践

1. **初始化时机**：在应用程序启动时尽早初始化日志系统
2. **defer Sync**：使用 `defer zlog.Sync()` 确保程序退出时日志正确落盘
3. **日志级别**：开发环境使用debug级别，生产环境使用info或warn级别
4. **结构化日志**：生产环境推荐使用结构化日志，便于日志分析工具处理
5. **日志字段**：添加足够的上下文信息，如用户ID、请求ID等，便于问题排查
6. **敏感信息**：避免在日志中记录密码、密钥等敏感信息

## 性能考虑

- 结构化日志（Debug/Info/Warn等）性能最高
- 键值对日志（Debugw/Infow等）性能次之
- 格式化日志（Debugf/Infof等）性能相对较低
- 对于高频日志，可以考虑启用日志采样功能

## 依赖

- [go.uber.org/zap](https://github.com/uber-go/zap)：High performance logging library
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack)：Log file rotation

## 注意事项

1. 确保在应用程序退出前调用`zlog.Sync()`来刷新所有日志到磁盘
2. 在生产环境中，建议将日志级别设置为`info`或更高，以减少日志量
3. 对于高频日志，考虑启用采样功能以提高性能

## 许可证

[MIT](https://opensource.org/licenses/MIT)