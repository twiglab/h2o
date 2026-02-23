package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// DeptRepository 部门数据仓库
type DeptRepository struct {
	db *gorm.DB
}

// NewDeptRepository 创建部门仓库
func NewDeptRepository(db *gorm.DB) *DeptRepository {
	return &DeptRepository{db: db}
}

// FindByID 根据ID查找部门
func (r *DeptRepository) FindByID(id int64) (*models.Dept, error) {
	var dept models.Dept
	if err := r.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

// FindByCode 根据编码查找部门
func (r *DeptRepository) FindByCode(code string) (*models.Dept, error) {
	var dept models.Dept
	if err := r.db.Where("dept_code = ?", code).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

// ExistsByCode 检查部门编码是否存在
func (r *DeptRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Dept{}).Where("dept_code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCodeExcludeID 检查部门编码是否存在（排除指定ID）
func (r *DeptRepository) ExistsByCodeExcludeID(code string, excludeID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Dept{}).Where("dept_code = ? AND id != ?", code, excludeID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 获取部门列表（不分页，用于构建树形结构）
func (r *DeptRepository) List(keyword string, status *int8) ([]models.Dept, error) {
	var depts []models.Dept

	query := r.db.Model(&models.Dept{}).Where("deleted_at IS NULL")

	if keyword != "" {
		query = query.Where("dept_name LIKE ? OR dept_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}

	return depts, nil
}

// FindChildren 查找子部门
func (r *DeptRepository) FindChildren(parentID int64) ([]models.Dept, error) {
	var depts []models.Dept
	if err := r.db.Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Order("sort_order ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

// HasChildren 检查是否有子部门
func (r *DeptRepository) HasChildren(id int64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Dept{}).Where("parent_id = ? AND deleted_at IS NULL", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create 创建部门
func (r *DeptRepository) Create(dept *models.Dept) error {
	return r.db.Create(dept).Error
}

// Update 更新部门
func (r *DeptRepository) Update(dept *models.Dept, updates map[string]interface{}) error {
	return r.db.Model(dept).Updates(updates).Error
}

// Delete 软删除部门
func (r *DeptRepository) Delete(id int64) error {
	return r.db.Model(&models.Dept{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
