package middleware

import (
	"errors"
	"strings"
	"time"

	"shared/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT Claims 结构(增强版)
type JWTClaims struct {
	UserID      int64    `json:"user_id"`
	Username    string   `json:"username"`
	RealName    string   `json:"real_name,omitempty"`
	UserType    int8     `json:"user_type"`
	DeptID      int64    `json:"dept_id,omitempty"`
	RoleIDs     []int64  `json:"role_ids,omitempty"`
	DataScope   int8     `json:"data_scope"`
	Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

// JWTMiddleware JWT 中间件
type JWTMiddleware struct {
	secret           []byte
	accessExpireMin  int
	refreshExpireDay int
}

// NewJWTMiddleware 创建 JWT 中间件
func NewJWTMiddleware(secret string, accessExpireMin, refreshExpireDay int) *JWTMiddleware {
	return &JWTMiddleware{
		secret:           []byte(secret),
		accessExpireMin:  accessExpireMin,
		refreshExpireDay: refreshExpireDay,
	}
}

// GenerateToken 生成 JWT token(增强版)
func (m *JWTMiddleware) GenerateToken(claims *JWTClaims) (accessToken, refreshToken string, err error) {
	now := time.Now()

	// Access Token
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.accessExpireMin) * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    "prepaid-utility",
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (只包含基本信息)
	refreshClaims := &JWTClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.refreshExpireDay) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "prepaid-utility-refresh",
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateTokenSimple 生成简单token(向后兼容)
func (m *JWTMiddleware) GenerateTokenSimple(userID int64, username string) (accessToken, refreshToken string, err error) {
	claims := &JWTClaims{
		UserID:    userID,
		Username:  username,
		UserType:  2,                 // 默认普通用户
		DataScope: DataScopeSelfOnly, // 默认仅本人
	}
	return m.GenerateToken(claims)
}

// ParseToken 解析 JWT token
func (m *JWTMiddleware) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AuthMiddleware 认证中间件(增强版)
func (m *JWTMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := m.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "认证失败")
			c.Abort()
			return
		}

		// 存储用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("real_name", claims.RealName)
		c.Set("user_type", claims.UserType)
		c.Set("dept_id", claims.DeptID)
		c.Set("role_ids", claims.RoleIDs)
		c.Set("data_scope", claims.DataScope)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件(不强制要求登录)
func (m *JWTMiddleware) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := m.ParseToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		// 存储用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("real_name", claims.RealName)
		c.Set("user_type", claims.UserType)
		c.Set("dept_id", claims.DeptID)
		c.Set("role_ids", claims.RoleIDs)
		c.Set("data_scope", claims.DataScope)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	if id, exists := c.Get("user_id"); exists {
		return id.(int64)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if name, exists := c.Get("username"); exists {
		return name.(string)
	}
	return ""
}

// GetRealName 从上下文获取真实姓名
func GetRealName(c *gin.Context) string {
	if name, exists := c.Get("real_name"); exists {
		return name.(string)
	}
	return ""
}

// GetClaims 从上下文获取完整Claims
func GetClaims(c *gin.Context) *JWTClaims {
	return &JWTClaims{
		UserID:      GetUserID(c),
		Username:    GetUsername(c),
		RealName:    GetRealName(c),
		UserType:    GetUserType(c),
		DeptID:      GetDeptID(c),
		RoleIDs:     GetRoleIDs(c),
		DataScope:   GetDataScope(c),
		Permissions: GetPermissions(c),
	}
}

// IsAuthenticated 检查是否已认证
func IsAuthenticated(c *gin.Context) bool {
	return GetUserID(c) > 0
}

// RefreshTokenResponse 刷新token响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // 秒
}

// RefreshToken 刷新token
func (m *JWTMiddleware) RefreshToken(refreshTokenString string, newClaims *JWTClaims) (*RefreshTokenResponse, error) {
	// 验证refresh token
	claims, err := m.ParseToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("refresh token无效")
	}

	// 验证issuer
	if claims.Issuer != "prepaid-utility-refresh" {
		return nil, errors.New("无效的refresh token")
	}

	// 使用新的claims生成token(如果提供),否则使用refresh token中的基本信息
	if newClaims == nil {
		newClaims = &JWTClaims{
			UserID:    claims.UserID,
			Username:  claims.Username,
			UserType:  2,
			DataScope: DataScopeSelfOnly,
		}
	}

	accessToken, refreshToken, err := m.GenerateToken(newClaims)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    m.accessExpireMin * 60,
	}, nil
}
