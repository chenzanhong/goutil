// Package jwtx provides JWT utilities integrated with Gin.
package jwtx

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SigningMethod = jwt.SigningMethod

// Signing methods aliases to avoid importing jwt package directly.
var (
	SigningMethodNone  SigningMethod = jwt.SigningMethodNone
	SigningMethodHS256 SigningMethod = jwt.SigningMethodHS256
	SigningMethodHS384 SigningMethod = jwt.SigningMethodHS384
	SigningMethodHS512 SigningMethod = jwt.SigningMethodHS512
	SigningMethodRS256 SigningMethod = jwt.SigningMethodRS256
	SigningMethodRS384 SigningMethod = jwt.SigningMethodRS384
	SigningMethodRS512 SigningMethod = jwt.SigningMethodRS512
	SigningMethodES256 SigningMethod = jwt.SigningMethodES256
	SigningMethodES384 SigningMethod = jwt.SigningMethodES384
	SigningMethodES512 SigningMethod = jwt.SigningMethodES512
	SigningMethodPS256 SigningMethod = jwt.SigningMethodPS256 // RSASSA-PSS
	SigningMethodPS384 SigningMethod = jwt.SigningMethodPS384
	SigningMethodPS512 SigningMethod = jwt.SigningMethodPS512
)

type ErrorType = error

var (
	ErrClaimsInvalid ErrorType = errors.New("claims must be a struct or pointer to struct")
	ErrSigningMethod ErrorType = errors.New("unexpected signing method")
	ErrInvalidKey              = jwt.ErrInvalidKey
	// ErrTokenNotValidYet                    = jwt.ErrTokenNotValidYet
	// ErrInvalidKeyType                      = jwt.ErrInvalidKeyType
	// ErrHashUnavailable                     = jwt.ErrHashUnavailable
	// ErrTokenMalformed                      = jwt.ErrTokenMalformed
	// ErrTokenUnverifiable                   = jwt.ErrTokenUnverifiable
	// ErrTokenSignatureInvalid               = jwt.ErrTokenSignatureInvalid
	// ErrTokenRequiredClaimMissing           = jwt.ErrTokenRequiredClaimMissing
	// ErrTokenInvalidAudience                = jwt.ErrTokenInvalidAudience
	// ErrTokenExpired                        = jwt.ErrTokenExpired
	// ErrTokenUsedBeforeIssued               = jwt.ErrTokenUsedBeforeIssued
	// ErrTokenInvalidIssuer                  = jwt.ErrTokenInvalidIssuer
	// ErrTokenInvalidSubject                 = jwt.ErrTokenInvalidSubject
	// ErrTokenInvalidId                      = jwt.ErrTokenInvalidId
	// ErrTokenInvalidClaims                  = jwt.ErrTokenInvalidClaims
	// ErrInvalidType                         = jwt.ErrInvalidType
)

// === 标准化的错误码常量 ===
const (
	ErrorCodeMissingToken   = "missing_token"
	ErrorCodeInvalidFormat  = "invalid_token_format"
	ErrorCodeTokenExpired   = "token_expired"
	ErrorCodeTokenNotActive = "token_not_active"
	ErrorCodeTokenMalformed = "token_malformed"
	ErrorCodeInvalidKey     = "invalid_key"
	ErrorCodeInvalidKeyType = "invalid_key_type"
	ErrorCodeTokenInvalid   = "token_invalid"
	ErrorCodeInternalError  = "internal_error"
)

// Claims is an example claims structure.
// Users should define their own claims with RegisteredClaims embedded.
//
// Example:
//
//	type MyClaims struct {
//	    UserID   uint   `json:"user_id"     inject:"user_id"`
//	    Username string `json:"username"`                     // will use "username" as context key
//	    Role     string                                           // no tag → uses field name "Role"
//	    jwtx.RegisteredClaims
//	}
//
// Then generate a token:
//
//	claims := &MyClaims{
//	    UserID: 123,
//	    Username: "alice",
//	    Role: "admin",
//	    RegisteredClaims: jwtx.RegisteredClaims{
//	        ExpiresAt: jwtx.NewNumericDate(time.Now().Add(24 * time.Hour)),
//	        IssuedAt:  jwtx.NewNumericDate(time.Now()),
//	        Issuer:    "myapp",
//	    },
//	}
//	token, err := jwtUtil.SignToken(claims)
type Claims = jwt.Claims
type RegisteredClaims = jwt.RegisteredClaims
type NumericDate = jwt.NumericDate
type ClaimStrings = jwt.ClaimStrings

// GinJWT holds configuration for JWT operations.
type GinJWT struct {
	JWTKey        []byte
	SigningMethod SigningMethod
	Claims        Claims        // Prototype for reflection; must be a pointer to a struct type.
	AutoInject    bool          // If true, automatically inject claim fields into gin.Context. Default: false.
	claimsFactory func() Claims // Factory function to create new claims instance.
}

// defaultGJWT is the package-level default instance.
// It is protected by a mutex to allow safe configuration before use.
var (
	defaultGJWT *GinJWT
	defaultInit sync.Once // ensure SetDefault* can only be called before first use
)

func NewNumericDate(t time.Time) *NumericDate {
	return jwt.NewNumericDate(t)
}

type Option func(*GinJWT)

func WithSigningMethod(method SigningMethod) Option {
	return func(g *GinJWT) {
		if method == nil {
			panic("jwtx: WithSigningMethod: method cannot be nil")
		}
		g.SigningMethod = method
	}
}

func buildClaimsFactory(claims Claims) func() Claims {
	claimsType := reflect.TypeOf(claims)
	if claimsType.Kind() == reflect.Ptr {
		claimsType = claimsType.Elem()
	}
	if claimsType.Kind() != reflect.Struct {
		panic("jwtx: WithClaims: claims must be a pointer to a struct type")
	}
	return func() Claims {
		return reflect.New(claimsType).Interface().(Claims)
	}
}

func WithAutoInject(enabled bool) Option {
	return func(g *GinJWT) {
		g.AutoInject = enabled
	}
}

// Init initializes the default GinJWT instance.
// Use SigningMethodHS256 as default.
// It should be called before any other jwtx function.
func Init(key string, method SigningMethod, claims Claims, opts ...Option) {
	g, err := NewGinJWT(key, method, claims, opts...)
	if err != nil {
		panic(err)
	}
	defaultInit.Do(func() {
		defaultGJWT = g
	})
}

// InitWithHS256 initializes the default instance using HS256.
// This is a convenience function for common use cases.
func InitWithHS256(key string, claims Claims, opts ...Option) {
	Init(key, SigningMethodHS256, claims, opts...)
}

func mustDefault() *GinJWT {
	if defaultGJWT == nil {
		panic("jwtx: default instance not initialized; call jwtx.Init(key, claims) first")
	}
	return defaultGJWT
}

// NewGinJWT creates a new GinJWT instance.
// The claims parameter should be a pointer to a zero-value struct (e.g., &MyClaims{}).
func NewGinJWT(key string, method SigningMethod, claims Claims, opts ...Option) (*GinJWT, error) {
	if key == "" {
		return nil, ErrInvalidKey
	}
	// 生产环境应该检查key长度是否符合要求
	// if strings.HasPrefix(defaultGJWT.SigningMethod.Alg(), "HS") {
	// 	minLen := map[string]int{"HS256": 32, "HS384": 48, "HS512": 64}[defaultGJWT.SigningMethod.Alg()]
	// 	if minLen > 0 && len(defaultGJWT.JWTKey) < minLen {
	// 		panic(fmt.Sprintf("jwtx: %s key must be at least %d bytes", defaultGJWT.SigningMethod.Alg(), minLen))
	// 	}
	// }
	if claims == nil {
		return nil, ErrClaimsInvalid
	}
	g := &GinJWT{
		JWTKey:        []byte(key),
		Claims:        claims,
		SigningMethod: method,
		AutoInject:    false,
	}
	g.claimsFactory = buildClaimsFactory(claims)

	// 应用用户选项
	for _, opt := range opts {
		opt(g)
	}
	return g, nil
}

// SignToken signs the given claims using the default configuration.
func SignToken(claims Claims) (string, error) {
	return mustDefault().SignToken(claims)
}

func GinJWTAuthMiddleware() gin.HandlerFunc {
	return mustDefault().GinJWTAuthMiddleware()
}

// ParseJWT parses a token string using the default configuration.
func ParseJWT(tokenStr string) (Claims, error) {
	return mustDefault().ParseJWT(tokenStr)
}

// SignToken generates a signed JWT string from the given claims.
func (g *GinJWT) SignToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(g.SigningMethod, claims)
	return token.SignedString(g.JWTKey)
}

// GinJWTAuthMiddleware returns a Gin middleware that validates JWT tokens.
// If AutoInject is enabled, it injects public claim fields into the context.
func (g *GinJWT) GinJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   ErrorCodeMissingToken,
				"message": "请求头中缺少 Token",
			})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   ErrorCodeTokenMalformed,
				"message": "Token 格式错误，缺少 Bearer 前缀",
			})
			c.Abort()
			return
		}

		tokenStr := authHeader[len(bearerPrefix):]

		// Create a new instance of the claims type via reflection.
		claimsType := reflect.TypeOf(g.Claims)
		if claimsType.Kind() == reflect.Ptr {
			claimsType = claimsType.Elem()
		}
		if claimsType.Kind() != reflect.Struct {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   ErrorCodeInternalError,
				"message": "Claims must be a struct or pointer to struct",
			})
			c.Abort()
			return
		}

		claims := g.claimsFactory()
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if t.Method != g.SigningMethod {
				return nil, ErrSigningMethod
			}
			return g.JWTKey, nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorCodeTokenExpired, "message": "Token 已过期"})
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorCodeTokenNotActive, "message": "Token 尚未激活"})
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorCodeTokenMalformed, "message": "Token 格式不正确"})
			case errors.Is(err, jwt.ErrInvalidKey):
				c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorCodeInvalidKey, "message": "无效的签名密钥"})
			case errors.Is(err, jwt.ErrInvalidKeyType):
				c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorCodeInvalidKeyType, "message": "签名密钥类型错误"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": ErrorCodeTokenInvalid, "message": "Token 无效: " + err.Error()})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   ErrorCodeTokenInvalid,
				"message": "无效的 Token",
			})
			c.Abort()
			return
		}

		if g.AutoInject {
			injectClaimsToContext(c, claims)
		}

		// Also store full claims in context for advanced usage.
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// injectClaimsToContext injects selected fields from claims into gin.Context.
// Priority for key name:
// 1. `inject:"custom_key"`
// 2. `json:"key"` (ignoring options like `,omitempty`)
// 3. Field name (e.g., "Role")
// Embedded fields (like RegisteredClaims) are skipped.
func injectClaimsToContext(c *gin.Context, claims Claims) {
	val := reflect.ValueOf(claims)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Skip embedded RegisteredClaims or other anonymous structs
		if fieldType.Anonymous {
			continue
		}

		var key string

		// 1. Try `inject` tag
		if injectTag := fieldType.Tag.Get("inject"); injectTag != "" {
			key = injectTag
		} else if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			// 2. Parse `json` tag (handle "name,omitempty")
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				key = jsonTag[:idx]
			} else {
				key = jsonTag
			}
			// Skip if json:"-"
			if key == "-" {
				continue
			}
		} else {
			// 3. Fallback to field name
			key = fieldType.Name
		}

		if key != "" {
			c.Set(key, field.Interface())
		}
	}
}

// ParseJWT parses a raw JWT string and returns the claims.
// Useful for non-middleware scenarios (e.g., WebSocket auth).
func (g *GinJWT) ParseJWT(tokenStr string) (Claims, error) {
	claimsType := reflect.TypeOf(g.Claims)
	if claimsType.Kind() == reflect.Ptr {
		claimsType = claimsType.Elem()
	}
	if claimsType.Kind() != reflect.Struct {
		return nil, ErrClaimsInvalid
	}

	claims := g.claimsFactory()

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != g.SigningMethod {
			return nil, ErrSigningMethod
		}
		return g.JWTKey, nil
	})

	if err != nil {
		return nil, err
	}
	// if !token.Valid {
	// 	return nil, ErrTokenInvalid
	// }
	_ = token
	return claims, nil
}

/*
	使用示例1，使用全局默认GinJWT
type MyClaims struct {
	UserID   uint   `json:"user_id" inject:"user_id"`
	Username string `json:"username"`
	Role     string
	jwtx.RegisteredClaims
}

func main() {
	claims := &MyClaims{}
	jwtx.Init(
		"my-32-byte-long-secret-key-1234567890ab",
		jwtx.SigningMethodHS256, // ← 必须传！
		claims,
		jwtx.WithAutoInject(true),
	)

	r := gin.Default()
	r.Use(jwtx.GinJWTAuthMiddleware())
	// ...
}
*/

/*
	使用示例2，自定义GinJWT
type MyClaims struct {
	UserID   uint   `json:"user_id" inject:"user_id"`
	Username string `json:"username"`
	Role     string
	jwtx.RegisteredClaims
}

func main() {
	claims := &MyClaims{}
	gjwt, err := jwtx.NewGinJWT("my-32-byte-long-secret-key-1234567890ab", jwt.SigningMethodHS256, claims)
	if err != nil {
		panic(err)
	}
	gjwt.AutoInject = true

	r := gin.Default()
	r.Use(gjwt.GinJWTAuthMiddleware())
	// ...
}

*/
