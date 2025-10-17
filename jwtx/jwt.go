package jwtx

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 默认配置 ----------------------------------------------
var (
	defaultJWTKey        = []byte(defaultKey)
	defaultSigningMethod = jwt.SigningMethodHS256
)

const defaultKey = "0123456789abcdef"

func GetDefaultSigningMethod() jwt.SigningMethod {
	return defaultSigningMethod
}

/*
Claims 示例结构体，原型。
推荐用户自定义Claims。如：

	claims := &model.MyClaims{
		UserID:   123,
		Username: "username",
		Role:     "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24小时有效期
			IssuedAt:  time.Now().Unix(),
			Issuer:    "myapp",
		},
	}

然后使用 SignToken 生成 token

	token, err := jwtUtil.SignToken(claims)

注意：有 jwt:"set"（默认以字段名为key）或 jwt:"set=custom_key"（以custom_key为key）的字段
才会 c.Set(key any, value any)
*/
type Claims struct {
	UserName string `json:"username" jwt:"set"`
	// …… 其他自定义字段
	jwt.StandardClaims
}

// GinJWT 结构体，支持自定义Claims、JWTKey 和 method ----------------------------------------------
type GinJWT struct {
	JWTKey        []byte
	SigningMethod jwt.SigningMethod
	Claims        jwt.Claims // 用户可以自定义Claims
	// Claims 是一个 claims 原型，用于反射创建新实例。
	// 它本身不存储任何运行时数据，不应被直接修改。
}

/*  用户完全可以自己定义
type MyClaims struct {
    UserID   uint   `json:"user_id" jwt:"set=user_id"`
    Role     string `json:"role" jwt:"set"`
    jwt.StandardClaims
}

jwtUtil := goutil.NewGinJWT("secret", jwt.SigningMethodHS256, MyClaims{})
*/

// newGinJWTWithClaims 创建一个新的GinJWT实例（内部使用）
func newGinJWTWithClaims(key []byte, method jwt.SigningMethod, claims jwt.Claims) *GinJWT {
	return &GinJWT{
		JWTKey:        key,
		SigningMethod: method,
		Claims:        claims,
	}
}

// NewGinJWT 创建一个新的GinJWT实例
func NewGinJWT(key string, method jwt.SigningMethod, claims jwt.Claims) *GinJWT {
	return newGinJWTWithClaims([]byte(key), method, claims)
}

// NewDefaultGinJWT 创建一个使用默认配置的新GinJWT实例
func NewDefaultGinJWT() *GinJWT {
	return newGinJWTWithClaims(defaultJWTKey, defaultSigningMethod, &Claims{})
}

// SignToken
func (g *GinJWT) SignToken(claims jwt.Claims) (string, error) {

	// 创建token对象
	token := jwt.NewWithClaims(g.SigningMethod, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString(g.JWTKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (g *GinJWT) GinJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_token", "message": "请求头中缺少Token"})
			c.Abort()
			return
		}

		// 处理 Bearer 前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(tokenStr, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token_format", "message": "Token格式错误，缺少Bearer前缀"})
			c.Abort()
			return
		}
		tokenStr = tokenStr[len(bearerPrefix):] // 去掉前缀

		// 使用反射创建g.Claims类型的新实例
		claims := reflect.New(reflect.TypeOf(g.Claims).Elem()).Interface().(jwt.Claims)

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// 关键安全检查：验证签名算法是否符合预期
			if token.Method != g.SigningMethod {
				return nil, errors.New("unexpected signing method")
			}
			return g.JWTKey, nil
		})

		// 详细错误处理
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				switch ve.Errors {
				case jwt.ValidationErrorExpired:
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token_expired", "message": "Token已过期"})
				case jwt.ValidationErrorNotValidYet:
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token_not_active", "message": "Token尚未激活"})
				case jwt.ValidationErrorMalformed:
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token_malformed", "message": "Token格式不正确"})
				default:
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token_invalid", "message": "Token无效: " + err.Error()})
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "parse_error", "message": "Token解析失败: " + err.Error()})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token_invalid", "message": "无效的Token"})
			c.Abort()
			return
		}

		extractClaimsToContext(c, claims)
		c.Next()
	}
}

// extractClaimsToContext 将 claims 中带有 jwt:"set" tag 的字段注入到 gin.Context
func extractClaimsToContext(c *gin.Context, claims jwt.Claims) {
	val := reflect.ValueOf(claims).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过非导出字段
		if !field.CanInterface() {
			continue
		}

		// 检查是否有jwt:"set" 或 jwt:"set=custom_key"
		jwtTag := fieldType.Tag.Get("jwt")
		if !strings.HasPrefix(jwtTag, "set") {
			continue
		}

		var key string
		if i := strings.Index(jwtTag, "="); i > 0 {
			key = jwtTag[i+1:]
		} else {
			key = fieldType.Name
		}

		c.Set(key, field.Interface())
	}
}

// ParseJWT 使用GinJWT解析JWT令牌
func (g *GinJWT) ParseJWT(tokenStr string) (jwt.Claims, error) {
	claims := reflect.New(reflect.TypeOf(g.Claims).Elem()).Interface().(jwt.Claims)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != g.SigningMethod {
			return nil, errors.New("unexpected signing method")
		}
		return g.JWTKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

/*
	使用示例：
	1、自定义Claims结构体
	type MyClaims struct {
		UserID   uint   `json:"user_id" jwt:"set=user_id"`     // 会自动 c.Set("user_id", value)
		Username string `json:"username" jwt:"set=username"`   // 会自动 c.Set("username", value)
		Role     string `json:"role" jwt:"set"`                // 会自动 c.Set("Role", value) 。用字段名做key
		jwt.StandardClaims
	}

	2. 初始化 JWT 工具
	var jwtUtil *goutil.GinJWT
	func init() {
		// 创建 JWT 工具，传入原型（注意：是 &model.MyClaims{}，不是实例数据）
		jwtUtil = goutil.NewGinJWT(
			"your-very-secret-key-32bytes-min!", // 建议至少 32 字节
			goutil.GetDefaultSigningMethod(),    // 或直接用 jwt.SigningMethodHS256
			&model.MyClaims{},                   // 传入结构体指针作为原型
		)
	}

	3. 登录接口 LoginHandler ：生成 Token
	// 创建 claims 实例（包含动态数据）
    claims := &model.MyClaims{
        UserID:   123,
        Username: req.Username,
        Role:     "admin",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24小时有效期
            IssuedAt:  time.Now().Unix(),
            Issuer:    "myapp",
        },
    }

    // 使用 SignToken 生成 token
    token, err := jwtUtil.SignToken(claims)

	4. 受保护的接口：使用中间件
	func ProfileHandler(c *gin.Context) {
		//  中间件已自动将 claims 中标记 jwt:"set" 的字段注入到 Context
		userID := c.GetUint("user_id")      // 来自 jwt:"set=user_id"
		username := c.GetString("username") // 来自 jwt:"set=username"
		role := c.GetString("Role")         // 来自 jwt:"set"

		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
			"role":     role,
		})
	}

	5. 路由设置
	// main.go
	func main() {
		r := gin.Default()

		// 公共路由
		r.POST("/login", LoginHandler)

		// 受保护的路由
		protected := r.Group("/api")
		protected.Use(jwtUtil.GinJWTAuthMiddleware()) // 使用中间件
		{
			protected.GET("/profile", ProfileHandler)
		}

		r.Run(":8080")
	}
*/
