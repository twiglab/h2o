package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataScope 数据权限范围常量
const (
	DataScopeAll       = 1 // 全部数据
	DataScopeDeptBelow = 2 // 本部门及下级部门
	DataScopeDeptOnly  = 3 // 仅本部门
	DataScopeSelfOnly  = 4 // 仅本人创建的数据
)

// DataScopeContext 数据权限上下文
type DataScopeContext struct {
	UserID    int64   `json:"user_id"`
	DeptID    int64   `json:"dept_id"`
	DataScope int8    `json:"data_scope"`
	DeptIDs   []int64 `json:"dept_ids"` // 有权限的部门ID列表(用于data_scope=2)
}

// GetDataScope 从上下文获取数据权限范围
func GetDataScope(c *gin.Context) int8 {
	if dataScope, exists := c.Get("data_scope"); exists {
		return dataScope.(int8)
	}
	return DataScopeSelfOnly
}

// GetDeptID 从上下文获取部门ID
func GetDeptID(c *gin.Context) int64 {
	if deptID, exists := c.Get("dept_id"); exists {
		return deptID.(int64)
	}
	return 0
}

// GetDataScopeContext 获取完整的数据权限上下文
func GetDataScopeContext(c *gin.Context) *DataScopeContext {
	ctx := &DataScopeContext{
		UserID:    GetUserID(c),
		DeptID:    GetDeptID(c),
		DataScope: GetDataScope(c),
	}

	// 获取部门ID列表(如果存在)
	if deptIDs, exists := c.Get("dept_ids"); exists {
		ctx.DeptIDs = deptIDs.([]int64)
	}

	return ctx
}

// DataScopeFilter 数据权限过滤器
// 用于在GORM查询中添加数据权限过滤条件
type DataScopeFilter struct {
	UserIDField   string // 用户ID字段名,默认"created_by"
	DeptIDField   string // 部门ID字段名,默认"dept_id"
	CustomerField string // 客户ID字段名,默认"customer_id"
}

// DefaultDataScopeFilter 默认数据权限过滤器
func DefaultDataScopeFilter() *DataScopeFilter {
	return &DataScopeFilter{
		UserIDField:   "created_by",
		DeptIDField:   "dept_id",
		CustomerField: "customer_id",
	}
}

// ApplyScope 应用数据权限过滤
func (f *DataScopeFilter) ApplyScope(c *gin.Context, db *gorm.DB) *gorm.DB {
	ctx := GetDataScopeContext(c)

	switch ctx.DataScope {
	case DataScopeAll:
		// 全部数据,不做过滤
		return db
	case DataScopeDeptBelow:
		// 本部门及下级部门
		if len(ctx.DeptIDs) > 0 {
			return db.Where(f.DeptIDField+" IN ?", ctx.DeptIDs)
		}
		return db.Where(f.DeptIDField+" = ?", ctx.DeptID)
	case DataScopeDeptOnly:
		// 仅本部门
		return db.Where(f.DeptIDField+" = ?", ctx.DeptID)
	case DataScopeSelfOnly:
		// 仅本人创建的数据
		return db.Where(f.UserIDField+" = ?", ctx.UserID)
	default:
		// 默认仅本人
		return db.Where(f.UserIDField+" = ?", ctx.UserID)
	}
}

// ApplyScopeByCustomer 按客户ID过滤数据权限
// 适用于三户模型中基于customer_id的数据过滤
func (f *DataScopeFilter) ApplyScopeByCustomer(c *gin.Context, db *gorm.DB, customerIDs []int64) *gorm.DB {
	ctx := GetDataScopeContext(c)

	switch ctx.DataScope {
	case DataScopeAll:
		// 全部数据,不做过滤
		return db
	default:
		// 其他情况按客户ID过滤
		if len(customerIDs) > 0 {
			return db.Where(f.CustomerField+" IN ?", customerIDs)
		}
		// 如果没有客户ID列表,返回空结果
		return db.Where("1 = 0")
	}
}

// ScopeFunc 返回GORM作用域函数
func (f *DataScopeFilter) ScopeFunc(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return f.ApplyScope(c, db)
	}
}

// DataScopeMiddleware 数据权限中间件
// 用于在请求中设置数据权限相关信息
func DataScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 数据权限信息已在JWT中间件中设置
		// 此中间件可用于额外的数据权限处理逻辑
		c.Next()
	}
}

// WithDataScope 便捷方法:在GORM查询中应用数据权限
func WithDataScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return DefaultDataScopeFilter().ScopeFunc(c)
}

// WithCustomDataScope 便捷方法:使用自定义字段的数据权限过滤
func WithCustomDataScope(c *gin.Context, userField, deptField string) func(db *gorm.DB) *gorm.DB {
	filter := &DataScopeFilter{
		UserIDField: userField,
		DeptIDField: deptField,
	}
	return filter.ScopeFunc(c)
}

// NeedDataFilter 检查是否需要数据权限过滤
func NeedDataFilter(c *gin.Context) bool {
	dataScope := GetDataScope(c)
	return dataScope != DataScopeAll
}

// CanAccessAllData 检查是否可以访问所有数据
func CanAccessAllData(c *gin.Context) bool {
	dataScope := GetDataScope(c)
	return dataScope == DataScopeAll
}

// CanAccessDeptData 检查是否可以访问部门数据
func CanAccessDeptData(c *gin.Context) bool {
	dataScope := GetDataScope(c)
	return dataScope == DataScopeAll || dataScope == DataScopeDeptBelow || dataScope == DataScopeDeptOnly
}
