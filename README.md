# goutil

一个 Go 工具库。


## 📦 当前功能

### ✅ JWT 认证助手（Gin + jwt-go）

为 Gin 框架量身打造的 JWT 工具，简化认证流程。

- 支持自定义 Claims 结构
- 通过 `jwt:"set"` 标签自动将字段注入 `gin.Context`
- 自动解析 `Bearer` 前缀
- 详细的 Token 错误处理（过期、格式错误等）
- 可自定义密钥和签名算法

#### 使用示例

```go
type MyClaims struct {
    UserID   uint   `json:"user_id" jwt:"set=user_id"` // 注入到 c.Set("user_id", ...)
    Role     string `json:"role" jwt:"set"`            // 注入到 c.Set("Role", ...)
    // 其他自定义字段
    jwt.StandardClaims // 该字段必须存在
}

// 初始化 JWT 工具
jwtUtil := goutil.NewGinJWT(
    "your-super-secret-key", 
    jwt.SigningMethodHS256, 
    &MyClaims{},
)

// 1. 生成 Token
claims := &MyClaims{
    UserID: 123,
    Role:   "admin",
    StandardClaims: jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24小时有效期
    },
}
token, err := jwtUtil.SignToken(claims)

// 2. 注册中间件
r.Use(jwtUtil.GinJWTAuthMiddleware())

// 3. 在 Handler 中使用
func ProfileHandler(c *gin.Context) {
    userID := c.GetUint("user_id") // 自动注入
    role := c.GetString("Role")    // 自动注入
    c.JSON(200, gin.H{"user_id": userID, "role": role})
}
```

---

### ✅ 动态路径解析器 `JoinPathFromCaller`

根据**调用文件的位置**，动态构建项目内文件的绝对路径。

告别 `cwd` 依赖和硬编码路径，轻松加载配置、模板、静态资源等。

#### 使用场景

你想从 `/your-project/handlers/user.go` 加载 `/your-project/config/db.yaml`。

#### 示例

```go
// 在 /your-project/handlers/user.go 中调用：
path, err := goutil.JoinPathFromCaller("..", "config", "db.yaml")
// 结果: /your-project/config/db.yaml
```

> ✅ 无论从哪个目录启动程序，路径都正确！

---

### ✅ 邮件发送工具

提供简单易用的邮件发送功能，支持HTML格式内容。

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

### ✅ 验证码生成工具

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

- 邮箱格式验证
- 手机号格式验证（待实现）

#### 使用示例

```go
// 验证邮箱格式
if regexx.IsValidEmail("test@example.com") {
    fmt.Println("邮箱格式正确")
} else {
    fmt.Println("邮箱格式不正确")
}
```

---

### ✅ SSH/SFTP 文件传输工具

提供基于SSH/SFTP的文件上传和下载功能。

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

### ✅ 高性能日志工具

基于Uber的zap日志库，提供高性能、结构化的日志记录功能。

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

// logx.RotateConfig{} 零值，表示不使用日志轮转功能
logger, _ := logx.SetupZapLogger("app.log", zapcore.InfoLevel, logx.RotateConfig{}, true, false)
```

---

### ✅ 定时任务工具

提供简单的定时任务执行功能。

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

## 🚀 快速开始

### 安装

```bash
go get github.com/chenzanhong/goutil
```

### 导入

```go
import "github.com/chenzanhong/goutil"
```
