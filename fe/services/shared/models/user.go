package models

import "time"

// 用户类型
const (
	UserTypeSuperAdmin = 1 // 超级管理员
	UserTypeNormal     = 2 // 普通用户
)

// 用户状态
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusNormal   = 1 // 正常
)

// User 系统用户(商管公司员工)
type User struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	Username    string     `gorm:"column:username" json:"username"`
	Password    string     `gorm:"column:password" json:"-"`
	RealName    *string    `gorm:"column:real_name" json:"real_name"`
	Phone       *string    `gorm:"column:phone" json:"phone"`
	Email       *string    `gorm:"column:email" json:"email"`
	Avatar      *string    `gorm:"column:avatar" json:"avatar"`
	DeptID      *int64     `gorm:"column:dept_id" json:"dept_id"`
	UserType    int8       `gorm:"column:user_type" json:"user_type"`
	Status      int8       `gorm:"column:status" json:"status"`
	LastLoginAt *time.Time `gorm:"column:last_login_at" json:"last_login_at"`
	LastLoginIP *string    `gorm:"column:last_login_ip" json:"last_login_ip"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Dept  *Dept  `gorm:"foreignKey:DeptID" json:"dept,omitempty"`
	Roles []Role `gorm:"many2many:sys_user_role;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:role_id" json:"roles,omitempty"`
}

func (User) TableName() string { return "sys_user" }

// IsSuperAdmin 是否为超级管理员
func (u *User) IsSuperAdmin() bool {
	return u.UserType == UserTypeSuperAdmin
}

// IsEnabled 是否启用
func (u *User) IsEnabled() bool {
	return u.Status == UserStatusNormal
}

// GetRoleIDs 获取用户所有角色ID
func (u *User) GetRoleIDs() []int64 {
	ids := make([]int64, 0, len(u.Roles))
	for _, role := range u.Roles {
		ids = append(ids, role.ID)
	}
	return ids
}

// GetPermissionCodes 获取用户所有权限码
func (u *User) GetPermissionCodes() []string {
	codes := make([]string, 0)
	seen := make(map[string]bool)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if !seen[perm.PermCode] {
				codes = append(codes, perm.PermCode)
				seen[perm.PermCode] = true
			}
		}
	}
	return codes
}

// GetDataScope 获取用户数据权限范围(取最大范围)
func (u *User) GetDataScope() int8 {
	if u.IsSuperAdmin() {
		return DataScopeAll
	}
	dataScope := int8(DataScopeSelfOnly) // 默认仅本人
	for _, role := range u.Roles {
		if role.DataScope < dataScope {
			dataScope = role.DataScope
		}
	}
	return dataScope
}

// HasPermission 是否拥有指定权限
func (u *User) HasPermission(permCode string) bool {
	if u.IsSuperAdmin() {
		return true
	}
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.PermCode == permCode {
				return true
			}
		}
	}
	return false
}

// UserLoginInfo 用户登录信息(用于JWT Claims)
type UserLoginInfo struct {
	UserID      int64    `json:"user_id"`
	Username    string   `json:"username"`
	RealName    string   `json:"real_name"`
	UserType    int8     `json:"user_type"`
	DeptID      int64    `json:"dept_id"`
	RoleIDs     []int64  `json:"role_ids"`
	DataScope   int8     `json:"data_scope"`
	Permissions []string `json:"permissions"`
}

// ToLoginInfo 转换为登录信息
func (u *User) ToLoginInfo() *UserLoginInfo {
	info := &UserLoginInfo{
		UserID:      u.ID,
		Username:    u.Username,
		UserType:    u.UserType,
		RoleIDs:     u.GetRoleIDs(),
		DataScope:   u.GetDataScope(),
		Permissions: u.GetPermissionCodes(),
	}
	if u.RealName != nil {
		info.RealName = *u.RealName
	}
	if u.DeptID != nil {
		info.DeptID = *u.DeptID
	}
	return info
}
