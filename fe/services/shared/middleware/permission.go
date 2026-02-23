package middleware

import (
	"strings"

	"shared/response"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 功能权限检查中间件
// permCode: 需要的权限编码,支持多个(逗号分隔,任一匹配即可)
func PermissionMiddleware(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户类型
		userType, exists := c.Get("user_type")
		if !exists {
			response.Forbidden(c, "无法获取用户信息")
			c.Abort()
			return
		}

		// 超级管理员跳过权限检查
		if userType.(int8) == 1 {
			c.Next()
			return
		}

		// 获取用户权限列表
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Forbidden(c, "无法获取权限信息")
			c.Abort()
			return
		}

		userPerms := permissions.([]string)
		requiredPerms := strings.Split(permCode, ",")

		// 检查是否拥有任一所需权限
		hasPermission := false
		for _, required := range requiredPerms {
			required = strings.TrimSpace(required)
			for _, userPerm := range userPerms {
				if userPerm == required {
					hasPermission = true
					break
				}
				// 支持通配符匹配: system:* 匹配 system:user:list
				if strings.HasSuffix(userPerm, ":*") {
					prefix := strings.TrimSuffix(userPerm, "*")
					if strings.HasPrefix(required, prefix) {
						hasPermission = true
						break
					}
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "没有操作权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissions 多权限检查(必须同时拥有所有权限)
func RequirePermissions(permCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户类型
		userType, exists := c.Get("user_type")
		if !exists {
			response.Forbidden(c, "无法获取用户信息")
			c.Abort()
			return
		}

		// 超级管理员跳过权限检查
		if userType.(int8) == 1 {
			c.Next()
			return
		}

		// 获取用户权限列表
		permissions, exists := c.Get("permissions")
		if !exists {
			response.Forbidden(c, "无法获取权限信息")
			c.Abort()
			return
		}

		userPerms := permissions.([]string)
		userPermMap := make(map[string]bool)
		for _, p := range userPerms {
			userPermMap[p] = true
		}

		// 检查是否拥有所有所需权限
		for _, required := range permCodes {
			if !userPermMap[required] {
				response.Forbidden(c, "没有操作权限")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// HasPermission 检查上下文中是否有指定权限
func HasPermission(c *gin.Context, permCode string) bool {
	// 获取用户类型
	userType, exists := c.Get("user_type")
	if !exists {
		return false
	}

	// 超级管理员拥有所有权限
	if userType.(int8) == 1 {
		return true
	}

	// 获取用户权限列表
	permissions, exists := c.Get("permissions")
	if !exists {
		return false
	}

	userPerms := permissions.([]string)
	for _, userPerm := range userPerms {
		if userPerm == permCode {
			return true
		}
		// 支持通配符匹配
		if strings.HasSuffix(userPerm, ":*") {
			prefix := strings.TrimSuffix(userPerm, "*")
			if strings.HasPrefix(permCode, prefix) {
				return true
			}
		}
	}

	return false
}

// GetUserType 获取用户类型
func GetUserType(c *gin.Context) int8 {
	if userType, exists := c.Get("user_type"); exists {
		return userType.(int8)
	}
	return 0
}

// IsSuperAdmin 是否为超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	return GetUserType(c) == 1
}

// GetRoleIDs 获取用户角色ID列表
func GetRoleIDs(c *gin.Context) []int64 {
	if roleIDs, exists := c.Get("role_ids"); exists {
		return roleIDs.([]int64)
	}
	return nil
}

// GetPermissions 获取用户权限列表
func GetPermissions(c *gin.Context) []string {
	if permissions, exists := c.Get("permissions"); exists {
		return permissions.([]string)
	}
	return nil
}
