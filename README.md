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

## 🚀 快速开始

### 安装

```bash
go get github.com/chenzanhong/goutil
```

### 导入

```go
import "github.com/chenzanhong/goutil"
```
