package models

import "time"

// 权限类型
const (
	PermTypeDirectory = 1 // 目录
	PermTypeMenu      = 2 // 菜单
	PermTypeButton    = 3 // 按钮/API
)

// 权限状态
const (
	PermStatusDisabled = 0 // 停用
	PermStatusNormal   = 1 // 正常
)

// 可见状态
const (
	PermVisibleNo  = 0 // 不可见
	PermVisibleYes = 1 // 可见
)

// Permission 权限表
type Permission struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	ParentID  int64     `gorm:"column:parent_id" json:"parent_id"`
	PermCode  string    `gorm:"column:perm_code" json:"perm_code"`
	PermName  string    `gorm:"column:perm_name" json:"perm_name"`
	PermType  int8      `gorm:"column:perm_type" json:"perm_type"`
	Path      *string   `gorm:"column:path" json:"path"`
	Component *string   `gorm:"column:component" json:"component"`
	Icon      *string   `gorm:"column:icon" json:"icon"`
	APIPath   *string   `gorm:"column:api_path" json:"api_path"`
	APIMethod *string   `gorm:"column:api_method" json:"api_method"`
	Visible   int8      `gorm:"column:visible" json:"visible"`
	SortOrder int       `gorm:"column:sort_order" json:"sort_order"`
	Status    int8      `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Parent   *Permission  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Permission `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Permission) TableName() string { return "sys_permission" }

// IsDirectory 是否为目录
func (p *Permission) IsDirectory() bool {
	return p.PermType == PermTypeDirectory
}

// IsMenu 是否为菜单
func (p *Permission) IsMenu() bool {
	return p.PermType == PermTypeMenu
}

// IsButton 是否为按钮/API
func (p *Permission) IsButton() bool {
	return p.PermType == PermTypeButton
}

// IsVisible 是否可见
func (p *Permission) IsVisible() bool {
	return p.Visible == PermVisibleYes
}

// IsEnabled 是否启用
func (p *Permission) IsEnabled() bool {
	return p.Status == PermStatusNormal
}

// IsRoot 是否为根权限
func (p *Permission) IsRoot() bool {
	return p.ParentID == 0
}

// PermissionTree 权限树结构
type PermissionTree struct {
	Permission
	Children []*PermissionTree `json:"children,omitempty"`
}

// MenuMeta 菜单元信息(用于前端路由)
type MenuMeta struct {
	Title string `json:"title"`
	Icon  string `json:"icon,omitempty"`
}

// MenuRoute 菜单路由(用于前端)
type MenuRoute struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Path      string       `json:"path"`
	Component string       `json:"component,omitempty"`
	Meta      MenuMeta     `json:"meta"`
	Children  []*MenuRoute `json:"children,omitempty"`
}
