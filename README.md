# Go Util - JWT 工具类

这是一个用于处理 JWT (JSON Web Tokens) 的 Go 工具包，专为 Gin 框架设计。

## 功能特性

- JWT 令牌生成
- JWT 令牌验证
- Gin 中间件集成
- 详细的错误处理
- 自定义声明支持
- 可扩展的 GinJWT 结构体

## 安装

```bash
go get github.com/dgrijalva/jwt-go
```

## 使用方法

### 1. 导入包

```go
import "your-module-path/go-util"
```

### 2. 设置 JWT 密钥

```go
// 设置自定义密钥
goutil.SetJWTKey("your-secret-key")
```

### 3. 生成 JWT 令牌

```go
// 生成 JWT 令牌
expirationTime := time.Now().Add(24 * time.Hour) // 24小时后过期
tokenString, err := goutil.GenerateJWT("username", 123, expirationTime)
if err != nil {
    // 处理错误
}
```

### 4. 使用 Gin 中间件

```go
// 在路由中使用中间件
r := gin.Default()
r.Use(goutil.GinJWTAuthMiddleware())

r.GET("/protected", func(c *gin.Context) {
    username := c.MustGet("username").(string)
    userID := c.MustGet("user_id").(uint)
    c.JSON(200, gin.H{
        "message": fmt.Sprintf("Hello %s, your user ID is %d", username, userID),
    })
})
```

### 5. 手动解析 JWT 令牌

```go
// 手动解析令牌
claims, err := goutil.ParseJWT(tokenString)
if err != nil {
    // 处理错误
}
fmt.Printf("Username: %s, UserID: %d\n", claims.Username, claims.UserID)
```

## 使用 GinJWT 结构体（推荐）

GinJWT 结构体提供了更高的可扩展性和自定义性：

### 1. 创建 GinJWT 实例

```go
// 创建默认的 GinJWT 实例
jwtHandler := goutil.NewGinJWT("your-secret-key")

// 或者创建带有自定义 Claims 的实例
customClaims := &MyCustomClaims{}
jwtHandler := goutil.NewGinJWTWithClaims("your-secret-key", customClaims)
```

### 2. 使用 GinJWT 生成令牌

```go
// 创建声明
claims := &goutil.Claims{
    Username: "username",
    UserID:   123,
    StandardClaims: jwt.StandardClaims{
        ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        Issuer:    "your-app",
    },
}

// 生成令牌
tokenString, err := jwtHandler.GenerateJWT(claims)
if err != nil {
    // 处理错误
}
```

### 3. 使用 GinJWT 中间件

```go
// 在路由中使用 GinJWT 中间件
r := gin.Default()
r.Use(jwtHandler.GinJWTAuthMiddleware())

r.GET("/protected", func(c *gin.Context) {
    claims := c.MustGet("claims").(jwt.Claims)
    c.JSON(200, gin.H{
        "message": "Access granted",
        "claims":  claims,
    })
})
```

### 4. 使用 GinJWT 解析令牌

```go
// 解析令牌
claims, err := jwtHandler.ParseJWT(tokenString)
if err != nil {
    // 处理错误
}
```

## 自定义 Claims

你可以创建自己的 Claims 结构体，只需要嵌入 `jwt.StandardClaims`：

```go
type MyCustomClaims struct {
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.StandardClaims
}

// 使用自定义 Claims
customClaims := &MyCustomClaims{
    Email: "user@example.com",
    Role:  "admin",
    StandardClaims: jwt.StandardClaims{
        ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        Issuer:    "your-app",
    },
}

jwtHandler := goutil.NewGinJWTWithClaims("your-secret-key", customClaims)
```

## 错误处理

中间件会根据不同的 JWT 错误类型返回相应的错误信息：
- `token_expired`: Token 已过期
- `token_not_active`: Token 尚未激活
- `token_malformed`: Token 格式不正确
- `token_invalid`: Token 无效
- `parse_error`: Token 解析失败
- `missing_token`: 请求头中缺少 Token

## 预定义声明

工具类包含以下预定义声明：
- `Username` (string): 用户名
- `UserID` (uint): 用户ID
- `StandardClaims`: 标准 JWT 声明

## 安全建议

1. 使用强密钥（推荐至少32个字符）
2. 设置合理的过期时间
3. 在 HTTPS 环境中传输令牌
4. 不要在 JWT 中存储敏感信息
5. 定期轮换密钥