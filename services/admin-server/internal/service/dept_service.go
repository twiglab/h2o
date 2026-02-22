package service

import (
	"errors"
	"time"

	"admin-server/internal/repository"
	"shared/models"
)

var (
	ErrDeptCodeExists      = errors.New("部门编码已存在")
	ErrDeptHasChildren     = errors.New("该部门下存在子部门，无法删除")
	ErrDeptNotFound        = errors.New("部门不存在")
	ErrCannotSetSelfParent = errors.New("不能将自己设置为父部门")
)

// DeptService 部门服务
type DeptService struct {
	deptRepo *repository.DeptRepository
}

// NewDeptService 创建部门服务
func NewDeptService(deptRepo *repository.DeptRepository) *DeptService {
	return &DeptService{deptRepo: deptRepo}
}

// DeptTreeNode 部门树节点
type DeptTreeNode struct {
	ID        int64           `json:"id"`
	ParentID  int64           `json:"parent_id"`
	DeptCode  string          `json:"dept_code"`
	DeptName  string          `json:"dept_name"`
	LeaderID  *int64          `json:"leader_id"`
	SortOrder int             `json:"sort_order"`
	Status    int8            `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	Children  []*DeptTreeNode `json:"children,omitempty"`
}

// List 获取部门列表（树形结构）
func (s *DeptService) List(keyword string, status *int8) ([]*DeptTreeNode, error) {
	depts, err := s.deptRepo.List(keyword, status)
	if err != nil {
		return nil, err
	}

	// 如果有搜索条件，返回扁平列表
	if keyword != "" {
		result := make([]*DeptTreeNode, len(depts))
		for i, dept := range depts {
			result[i] = s.deptToTreeNode(&dept)
		}
		return result, nil
	}

	// 构建树形结构
	return s.buildTree(depts), nil
}

// buildTree 构建树形结构
func (s *DeptService) buildTree(depts []models.Dept) []*DeptTreeNode {
	nodeMap := make(map[int64]*DeptTreeNode)
	var roots []*DeptTreeNode

	// 先创建所有节点
	for i := range depts {
		node := s.deptToTreeNode(&depts[i])
		nodeMap[node.ID] = node
	}

	// 建立父子关系
	for i := range depts {
		node := nodeMap[depts[i].ID]
		if depts[i].ParentID == 0 {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[depts[i].ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}

func (s *DeptService) deptToTreeNode(dept *models.Dept) *DeptTreeNode {
	return &DeptTreeNode{
		ID:        dept.ID,
		ParentID:  dept.ParentID,
		DeptCode:  dept.DeptCode,
		DeptName:  dept.DeptName,
		LeaderID:  dept.LeaderID,
		SortOrder: dept.SortOrder,
		Status:    dept.Status,
		CreatedAt: dept.CreatedAt,
		Children:  []*DeptTreeNode{},
	}
}

// GetByID 获取部门详情
func (s *DeptService) GetByID(id int64) (*models.Dept, error) {
	return s.deptRepo.FindByID(id)
}

// CreateDeptInput 创建部门输入
type CreateDeptInput struct {
	ParentID  int64
	DeptCode  string
	DeptName  string
	LeaderID  *int64
	SortOrder int
	Status    int8
}

// Create 创建部门
func (s *DeptService) Create(input *CreateDeptInput) (*models.Dept, error) {
	// 检查编码是否存在
	exists, err := s.deptRepo.ExistsByCode(input.DeptCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDeptCodeExists
	}

	// 检查父部门是否存在
	if input.ParentID > 0 {
		_, err := s.deptRepo.FindByID(input.ParentID)
		if err != nil {
			return nil, ErrDeptNotFound
		}
	}

	status := input.Status
	if status == 0 {
		status = models.DeptStatusNormal
	}

	dept := &models.Dept{
		ParentID:  input.ParentID,
		DeptCode:  input.DeptCode,
		DeptName:  input.DeptName,
		LeaderID:  input.LeaderID,
		SortOrder: input.SortOrder,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.deptRepo.Create(dept); err != nil {
		return nil, err
	}

	return dept, nil
}

// UpdateDeptInput 更新部门输入
type UpdateDeptInput struct {
	ParentID  *int64
	DeptCode  *string
	DeptName  *string
	LeaderID  *int64
	SortOrder *int
	Status    *int8
}

// Update 更新部门
func (s *DeptService) Update(id int64, input *UpdateDeptInput) (*models.Dept, error) {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		return nil, ErrDeptNotFound
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.ParentID != nil {
		// 不能将自己设为父部门
		if *input.ParentID == id {
			return nil, ErrCannotSetSelfParent
		}
		updates["parent_id"] = *input.ParentID
	}

	if input.DeptCode != nil {
		// 检查编码是否被其他部门使用
		exists, err := s.deptRepo.ExistsByCodeExcludeID(*input.DeptCode, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrDeptCodeExists
		}
		updates["dept_code"] = *input.DeptCode
	}

	if input.DeptName != nil {
		updates["dept_name"] = *input.DeptName
	}
	if input.LeaderID != nil {
		updates["leader_id"] = *input.LeaderID
	}
	if input.SortOrder != nil {
		updates["sort_order"] = *input.SortOrder
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if err := s.deptRepo.Update(dept, updates); err != nil {
		return nil, err
	}

	return s.deptRepo.FindByID(id)
}

// Delete 删除部门
func (s *DeptService) Delete(id int64) error {
	// 检查是否有子部门
	hasChildren, err := s.deptRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return ErrDeptHasChildren
	}

	return s.deptRepo.Delete(id)
}
