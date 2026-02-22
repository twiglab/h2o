package service

import (
	"errors"
	"time"

	"admin-server/internal/repository"
	"shared/models"
)

var (
	ErrPermCodeExists  = errors.New("权限编码已存在")
	ErrPermNotFound    = errors.New("权限不存在")
	ErrPermHasChildren = errors.New("该权限下有子权限，无法删除")
)

// PermissionTree 权限树节点
type PermissionTree struct {
	ID        int64             `json:"id"`
	ParentID  int64             `json:"parent_id"`
	PermCode  string            `json:"perm_code"`
	PermName  string            `json:"perm_name"`
	PermType  int8              `json:"perm_type"`
	Path      *string           `json:"path"`
	Component *string           `json:"component"`
	Icon      *string           `json:"icon"`
	APIPath   *string           `json:"api_path"`
	APIMethod *string           `json:"api_method"`
	Visible   int8              `json:"visible"`
	SortOrder int               `json:"sort_order"`
	Status    int8              `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Children  []*PermissionTree `json:"children,omitempty"`
}

// PermissionService 权限服务
type PermissionService struct {
	permRepo *repository.PermissionRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService(permRepo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{permRepo: permRepo}
}

// GetTree 获取权限树
func (s *PermissionService) GetTree(status *int8) ([]*PermissionTree, error) {
	permissions, err := s.permRepo.List(status)
	if err != nil {
		return nil, err
	}

	return s.buildTree(permissions), nil
}

// buildTree 构建权限树
func (s *PermissionService) buildTree(permissions []models.Permission) []*PermissionTree {
	// 转换为树节点map
	nodeMap := make(map[int64]*PermissionTree)
	for _, p := range permissions {
		nodeMap[p.ID] = &PermissionTree{
			ID:        p.ID,
			ParentID:  p.ParentID,
			PermCode:  p.PermCode,
			PermName:  p.PermName,
			PermType:  p.PermType,
			Path:      p.Path,
			Component: p.Component,
			Icon:      p.Icon,
			APIPath:   p.APIPath,
			APIMethod: p.APIMethod,
			Visible:   p.Visible,
			SortOrder: p.SortOrder,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			Children:  make([]*PermissionTree, 0),
		}
	}

	// 构建树结构
	var roots []*PermissionTree
	for _, node := range nodeMap {
		if node.ParentID == 0 {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}

// GetByID 获取权限详情
func (s *PermissionService) GetByID(id int64) (*models.Permission, error) {
	return s.permRepo.FindByID(id)
}

// CreatePermissionInput 创建权限输入
type CreatePermissionInput struct {
	ParentID  int64
	PermCode  string
	PermName  string
	PermType  int8
	Path      *string
	Component *string
	Icon      *string
	APIPath   *string
	APIMethod *string
	Visible   int8
	SortOrder int
	Status    int8
}

// Create 创建权限
func (s *PermissionService) Create(input *CreatePermissionInput) (*models.Permission, error) {
	exists, err := s.permRepo.ExistsByCode(input.PermCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPermCodeExists
	}

	// 默认值
	visible := input.Visible
	if visible == 0 {
		visible = models.PermVisibleYes
	}
	status := input.Status
	if status == 0 {
		status = models.PermStatusNormal
	}

	permission := &models.Permission{
		ParentID:  input.ParentID,
		PermCode:  input.PermCode,
		PermName:  input.PermName,
		PermType:  input.PermType,
		Path:      input.Path,
		Component: input.Component,
		Icon:      input.Icon,
		APIPath:   input.APIPath,
		APIMethod: input.APIMethod,
		Visible:   visible,
		SortOrder: input.SortOrder,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.permRepo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// UpdatePermissionInput 更新权限输入
type UpdatePermissionInput struct {
	ParentID  *int64
	PermName  *string
	PermType  *int8
	Path      *string
	Component *string
	Icon      *string
	APIPath   *string
	APIMethod *string
	Visible   *int8
	SortOrder *int
	Status    *int8
}

// Update 更新权限
func (s *PermissionService) Update(id int64, input *UpdatePermissionInput) (*models.Permission, error) {
	permission, err := s.permRepo.FindByID(id)
	if err != nil {
		return nil, ErrPermNotFound
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.ParentID != nil {
		updates["parent_id"] = *input.ParentID
	}
	if input.PermName != nil {
		updates["perm_name"] = *input.PermName
	}
	if input.PermType != nil {
		updates["perm_type"] = *input.PermType
	}
	if input.Path != nil {
		updates["path"] = *input.Path
	}
	if input.Component != nil {
		updates["component"] = *input.Component
	}
	if input.Icon != nil {
		updates["icon"] = *input.Icon
	}
	if input.APIPath != nil {
		updates["api_path"] = *input.APIPath
	}
	if input.APIMethod != nil {
		updates["api_method"] = *input.APIMethod
	}
	if input.Visible != nil {
		updates["visible"] = *input.Visible
	}
	if input.SortOrder != nil {
		updates["sort_order"] = *input.SortOrder
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if err := s.permRepo.Update(permission, updates); err != nil {
		return nil, err
	}

	return s.permRepo.FindByID(id)
}

// Delete 删除权限
func (s *PermissionService) Delete(id int64) error {
	// 检查是否有子权限
	hasChildren, err := s.permRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return ErrPermHasChildren
	}

	return s.permRepo.Delete(id)
}
