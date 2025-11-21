# goutil

一个 Go 语言工具库。

## 📦 功能模块

### 🔐 JWT 认证助手（Gin + golang-jwt）

为 Gin 框架量身打造的 JWT 工具，简化认证流程。

#### 特性
- 支持自定义 Claims 结构
- 通过 `inject:"key"` 标签自动将字段注入 `gin.Context`
- 自动解析 `Bearer` 前缀
- 详细的 Token 错误处理（过期、格式错误等）
- 可自定义密钥和签名算法
- 支持自动注入声明字段到 Gin 上下文

#### 使用示例

```go
type MyClaims struct {
    UserID   uint   `json:"user_id" inject:"user_id"` // 注入到 c.Set("user_id", ...)
    Username string `json:"username"`                 // 将使用 "username" 作为上下文键
    Role     string                                   // 无标签 → 使用字段名 "Role"
    jwtx.RegisteredClaims
}

// 初始化 JWT 工具
jwtx.InitWithHS256("your-super-secret-key", &MyClaims{})

// 或者使用更多选项
jwtx.Init(
    "your-super-secret-key", 
    jwtx.SigningMethodHS256, 
    &MyClaims{},
    jwtx.WithAutoInject(true), // 启用自动注入
)

// 1. 生成 Token
claims := &MyClaims{
    UserID: 123,
    Username: "alice",
    Role: "admin",
    RegisteredClaims: jwtx.RegisteredClaims{
        ExpiresAt: jwtx.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24小时有效期
        IssuedAt:  jwtx.NewNumericDate(time.Now()),
        Issuer:    "myapp",
    },
}

token, err := jwtx.SignToken(claims)
if err != nil {
    // 处理错误
}

// 2. 注册中间件
r.Use(jwtx.GinJWTAuthMiddleware())

// 3. 在 Handler 中使用
func ProfileHandler(c *gin.Context) {
    userID := c.GetUint("user_id")     // 自动注入
    username := c.GetString("username") // 自动注入
    role := c.GetString("Role")        // 自动注入
    c.JSON(200, gin.H{"user_id": userID, "username": username, "role": role})
}
```

---

### 📍 动态路径解析器 `JoinPathFromCaller`

根据**调用文件的位置**，动态构建项目内文件的绝对路径。

告别 `cwd` 依赖和硬编码路径，轻松加载配置、模板、静态资源等。

#### 使用场景

你想从 `/your-project/handlers/user.go` 加载 `/your-project/config/db.yaml`。

#### 示例

```go
// 在 /your-project/handlers/user.go 中调用：
path, err := config.JoinPathFromCaller("..", "config", "db.yaml")
// 结果: /your-project/config/db.yaml
```


---

### 📧 邮件发送工具

提供简单易用的邮件发送功能，支持HTML格式内容。

#### 特性
- 支持常见的SMTP服务器
- 自动处理TLS连接
- 提供详细的错误信息反馈
- 支持自定义SMTP服务器地址和端口

#### 使用示例

```go
err := email.SendEmail(
    "your-email@example.com",     // 发送方邮箱
    "your-password",              // 发送方密码或授权码
    "recipient@example.com",      // 接收方邮箱
    "smtp.example.com",           // SMTP服务器地址
    587,                          // SMTP服务器端口
    "邮件主题",                   // 邮件主题
    "<h1>这是一封HTML邮件</h1>",  // 邮件内容（HTML格式）
)

if err != nil {
    log.Printf("邮件发送失败: %v", err)
} else {
    log.Println("邮件发送成功")
}
```

---

### 🔢 数学工具

提供常用的数学计算函数。

#### 特性
- 快速幂运算
- 组合数计算（支持模运算）
- 最大值/最小值比较
- 归并排序算法

#### 使用示例

```go
// 快速幂运算
result := mathx.PowMod(2, 10, 1000) // (2^10) % 1000

// 组合数计算
comb := mathx.Combination(10, 3, 1000000007) // C(10,3) % 1000000007

// 最大值/最小值
maxVal := mathx.Max(10, 20)
minVal := mathx.Min(10, 20)

// 归并排序
arr := []int{5, 2, 8, 1, 9}
sortedArr := mathx.MergeSort(arr)
```

---

### 🔐 验证码生成工具

生成安全的随机字符串，可用于验证码、临时密码等场景。

#### 特性
- 生成指定长度和字符集的随机字符串
- 支持自定义种子的随机字符串生成
- 可用于生成验证码、随机密码等

#### 使用示例

```go
// 生成6位数字验证码
code := randx.GenerateRandomToken(6, "0123456789")

// 生成包含字母和数字的8位随机密码
password := randx.GenerateRandomToken(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// 使用自定义种子生成随机字符串
token := randx.GenerateRandomTokenWithSeed(10, "abcdefghijklmnopqrstuvwxyz", time.Now().UnixNano())
```

---

### ✅ 正则表达式验证工具

提供常用的正则表达式验证功能，简化数据校验。

#### 特性
- 邮箱格式验证
- 手机号格式验证

#### 使用示例

```go
// 验证邮箱格式
if regexx.IsValidEmail("test@example.com") {
    fmt.Println("邮箱格式正确")
} else {
    fmt.Println("邮箱格式不正确")
}

// 验证手机号格式
if regexx.IsValidPhone("13812345678") {
    fmt.Println("手机号格式正确")
} else {
    fmt.Println("手机号格式不正确")
}
```

---

### 📁 SSH/SFTP 文件传输工具

提供基于SSH/SFTP的文件上传和下载功能。

#### 特性
- 支持文件上传到远程服务器
- 支持从远程服务器下载文件
- 简单易用的API接口

#### 使用示例

```go
// 上传文件到远程服务器
err := sshx.UploadFile(
    "./local/file.txt",           // 本地文件路径
    "/remote/path/file.txt",      // 远程文件路径
    "192.168.1.100:22",           // 服务器地址和端口
    "username",                   // 用户名
    "password",                   // 密码
)

if err != nil {
    log.Printf("文件上传失败: %v", err)
} else {
    log.Println("文件上传成功")
}

// 从远程服务器下载文件
err = sshx.DownloadFile(
    "/remote/path/file.txt",      // 远程文件路径
    "./local/downloaded.txt",     // 本地保存路径
    "192.168.1.100:22",           // 服务器地址和端口
    "username",                   // 用户名
    "password",                   // 密码
)

if err != nil {
    log.Printf("文件下载失败: %v", err)
} else {
    log.Println("文件下载成功")
}
```

---

### 🪵 日志工具 (zlog)

基于 Uber 的 zap 日志库，提供高性能、结构化的日志记录功能。

#### 特性
- 高性能：基于 uber-go/zap 和 lumberjack.v2 构建
- 多种日志风格：支持结构化日志、键值对日志和格式化日志
- 灵活配置：支持控制台输出、文件输出或同时输出
- 自动轮转：支持日志文件自动轮转和压缩
- 环境变量配置：支持通过环境变量进行配置
- 日志钩子：支持自定义日志钩子进行扩展
- 线程安全：全局实例的初始化是线程安全的

#### 使用示例

```go
// 默认初始化
err := zlog.InitLoggerDefault()
if err != nil {
    panic("初始化日志失败: " + err.Error())
}
defer zlog.Sync() // 确保日志正确刷盘

// 使用日志
zlog.Info("应用启动", zlog.String("version", "1.0.0"))
zlog.Infow("用户登录", "username", "admin", "ip", "127.0.0.1")
zlog.Infof("处理请求耗时: %v", 100*time.Millisecond)

// 自定义配置初始化
config := &zlog.LoggerConfig{
    Level:      "debug",          // 日志级别
    Output:     "both",           // 输出目标：console, file, both
    Format:     "json",           // 控制台格式：json, console
    FilePath:   "./logs/app.log", // 日志文件路径
    MaxSize:    100,              // 单个日志文件最大大小(MB)
    MaxBackups: 10,               // 保留的最大日志文件数
    MaxAge:     30,               // 保留的最大天数
    Compress:   true,             // 是否压缩旧日志文件
    Sampling:   false,            // 是否启用日志采样
}

err = zlog.InitLogger(config)
if err != nil {
    panic("初始化日志失败: " + err.Error())
}
defer zlog.Sync()

// 使用不同类型的日志
zlog.Debug("调试信息", zlog.Int("count", 10))
zlog.Infow("用户操作", "user", "admin", "action", "create", "id", 100)
zlog.Errorf("连接数据库失败: %v", err)
```

---

### 🪵 简单日志工具 (logx)

提供基于 zap 的简单日志配置工具。

#### 特性
- 支持日志轮转（自动分割和清理旧日志）
- JSON格式输出，便于日志收集和分析
- 可配置的日志级别
- 支持调用者信息和堆栈跟踪

#### 使用示例

```go
// 使用默认配置创建日志记录器
logger, err := logx.SetupDefaultZapLogger("./logs/app.log")
if err != nil {
    panic(err)
}
defer logger.Sync()

// 记录不同级别的日志
logger.Info("这是一条信息日志", zap.String("user", "张三"), zap.Int("age", 25))
logger.Warn("这是一条警告日志", zap.String("module", "auth"))
logger.Error("这是一条错误日志", zap.Error(errors.New("数据库连接失败")))

// 自定义配置创建日志记录器
customLogger, err := logx.SetupZapLogger(
    "./logs/custom.log",           // 日志文件路径
    zapcore.InfoLevel,             // 日志级别
    logx.RotateConfig{             // 轮转配置
        MaxSize:    20,            // 每个日志文件最大20MB
        MaxBackups: 10,            // 最多保留10个备份文件
        MaxAge:     30,            // 日志文件最多保留30天
        Compress:   true,          // 压缩旧日志文件
    },
    true,  // 添加调用者信息
    true,  // 添加堆栈跟踪
)
if err != nil {
    panic(err)
}
defer customLogger.Sync()

customLogger.Info("自定义配置日志", zap.String("service", "payment"))
```

---

### ⏱️ 定时任务工具

提供简单的定时任务执行功能。

#### 特性
- 支持以固定时间间隔执行函数
- 支持通过闭包传递参数

#### 使用示例

```go
// 每秒执行一次的定时任务
count := 0
go timex.Every(1*time.Second, func() {
    count++
    fmt.Printf("第 %d 个定时任务执行\n", count)
})

// 每分钟执行一次的任务
go timex.Every(1*time.Minute, func() {
    fmt.Println("每分钟执行的任务")
})
```

---

### 🔍 字符串搜索工具

提供高效的字符串搜索算法。

#### 特性
- KMP 字符串搜索算法
- 返回所有匹配位置的起始索引

#### 使用示例

```go
text := "ABABAABABA"
pattern := "ABA"

// 使用KMP算法搜索
positions := goutil.KMP(text, pattern)
fmt.Println("匹配位置:", positions) // 输出: [0 2 5 7]
```

---

## 🚀 快速开始

### 安装

```bash
go get github.com/chenzanhong/goutil
```

### 导入

```go
import "github.com/chenzanhong/goutil"
```

根据不同模块，也可以单独导入：

```go
import (
    "github.com/chenzanhong/goutil/config"
    "github.com/chenzanhong/goutil/email"
    "github.com/chenzanhong/goutil/jwtx"
    "github.com/chenzanhong/goutil/logx"
    "github.com/chenzanhong/goutil/mathx"
    "github.com/chenzanhong/goutil/randx"
    "github.com/chenzanhong/goutil/regexx"
    "github.com/chenzanhong/goutil/sshx"
    "github.com/chenzanhong/goutil/strings"
    "github.com/chenzanhong/goutil/timex"
    "github.com/chenzanhong/goutil/zlog"
)
```

## 📄 许可证

本项目采用 MIT 许可证，详情请见 [LICENSE](LICENSE) 文件。