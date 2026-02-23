package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// PermissionRepository 权限数据仓库
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// FindByID 根据ID查找权限
func (r *PermissionRepository) FindByID(id int64) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByCode 根据权限编码查找
func (r *PermissionRepository) FindByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.Where("perm_code = ?", code).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// ExistsByCode 检查权限编码是否存在
func (r *PermissionRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Permission{}).Where("perm_code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 获取所有权限（扁平列表）
func (r *PermissionRepository) List(status *int8) ([]models.Permission, error) {
	var permissions []models.Permission

	query := r.db.Model(&models.Permission{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// ListByParentID 获取指定父级下的权限
func (r *PermissionRepository) ListByParentID(parentID int64) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.Where("parent_id = ?", parentID).Order("sort_order ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// ListByIDs 根据ID列表获取权限
func (r *PermissionRepository) ListByIDs(ids []int64) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.Where("id IN ?", ids).Order("sort_order ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// ListMenus 获取菜单权限（目录和菜单类型）
func (r *PermissionRepository) ListMenus() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.Where("perm_type IN ? AND status = ? AND visible = ?",
		[]int8{models.PermTypeDirectory, models.PermTypeMenu},
		models.PermStatusNormal,
		models.PermVisibleYes,
	).Order("sort_order ASC, id ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// Create 创建权限
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// Update 更新权限
func (r *PermissionRepository) Update(permission *models.Permission, updates map[string]interface{}) error {
	return r.db.Model(permission).Updates(updates).Error
}

// Delete 删除权限
func (r *PermissionRepository) Delete(id int64) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

// HasChildren 检查是否有子权限
func (r *PermissionRepository) HasChildren(id int64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Permission{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// DB 返回数据库连接
func (r *PermissionRepository) DB() *gorm.DB {
	return r.db
}
