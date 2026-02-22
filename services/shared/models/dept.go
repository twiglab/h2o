package models

import "time"

// 部门状态
const (
	DeptStatusDisabled = 0 // 停用
	DeptStatusNormal   = 1 // 正常
)

// Dept 部门表
type Dept struct {
	ID        int64      `gorm:"column:id;primaryKey" json:"id"`
	ParentID  int64      `gorm:"column:parent_id" json:"parent_id"`
	DeptCode  string     `gorm:"column:dept_code" json:"dept_code"`
	DeptName  string     `gorm:"column:dept_name" json:"dept_name"`
	LeaderID  *int64     `gorm:"column:leader_id" json:"leader_id"`
	SortOrder int        `gorm:"column:sort_order" json:"sort_order"`
	Status    int8       `gorm:"column:status" json:"status"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Parent   *Dept  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Dept `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Leader   *User  `gorm:"foreignKey:LeaderID" json:"leader,omitempty"`
}

func (Dept) TableName() string { return "sys_dept" }

// IsRoot 是否为根部门
func (d *Dept) IsRoot() bool {
	return d.ParentID == 0
}

// IsEnabled 是否启用
func (d *Dept) IsEnabled() bool {
	return d.Status == DeptStatusNormal
}

// DeptTree 部门树结构
type DeptTree struct {
	Dept
	Children []*DeptTree `json:"children,omitempty"`
}
