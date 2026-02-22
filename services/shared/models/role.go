package models

import "time"

// 角色类型
const (
	RoleTypeBuiltin = 1 // 内置
	RoleTypeCustom  = 2 // 自定义
)

// 数据权限范围
const (
	DataScopeAll       = 1 // 全部数据
	DataScopeDeptBelow = 2 // 本部门及下级部门
	DataScopeDeptOnly  = 3 // 仅本部门
	DataScopeSelfOnly  = 4 // 仅本人创建的数据
)

// 角色状态
const (
	RoleStatusDisabled = 0 // 停用
	RoleStatusNormal   = 1 // 正常
)

// Role 角色表
type Role struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	RoleCode    string     `gorm:"column:role_code" json:"role_code"`
	RoleName    string     `gorm:"column:role_name" json:"role_name"`
	RoleType    int8       `gorm:"column:role_type" json:"role_type"`
	DataScope   int8       `gorm:"column:data_scope" json:"data_scope"`
	Description *string    `gorm:"column:description" json:"description"`
	Status      int8       `gorm:"column:status" json:"status"`
	SortOrder   int        `gorm:"column:sort_order" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Permissions []Permission `gorm:"many2many:sys_role_permission;foreignKey:ID;joinForeignKey:role_id;References:ID;joinReferences:permission_id" json:"permissions,omitempty"`
}

func (Role) TableName() string { return "sys_role" }

// IsBuiltin 是否为内置角色
func (r *Role) IsBuiltin() bool {
	return r.RoleType == RoleTypeBuiltin
}

// IsEnabled 是否启用
func (r *Role) IsEnabled() bool {
	return r.Status == RoleStatusNormal
}

// IsSuperAdmin 是否为超级管理员
func (r *Role) IsSuperAdmin() bool {
	return r.RoleCode == "super_admin"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	UserID    int64     `gorm:"column:user_id" json:"user_id"`
	RoleID    int64     `gorm:"column:role_id" json:"role_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (UserRole) TableName() string { return "sys_user_role" }

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           int64     `gorm:"column:id;primaryKey" json:"id"`
	RoleID       int64     `gorm:"column:role_id" json:"role_id"`
	PermissionID int64     `gorm:"column:permission_id" json:"permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (RolePermission) TableName() string { return "sys_role_permission" }
