package service

import (
	"errors"
	"time"

	"admin-server/internal/repository"
	"shared/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists    = errors.New("用户名已存在")
	ErrCannotDeleteAdmin = errors.New("不能删除管理员账号")
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// List 获取用户列表
func (s *UserService) List(keyword string, status *int8, page, pageSize int) ([]models.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(keyword, status, offset, pageSize)
}

// GetByID 获取用户详情
func (s *UserService) GetByID(id int64) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// CreateUserInput 创建用户输入
type CreateUserInput struct {
	Username string
	Password string
	RealName *string
	Phone    *string
	Email    *string
	DeptID   *int64
	UserType int8
	Status   int8
}

// Create 创建用户
func (s *UserService) Create(input *CreateUserInput) (*models.User, error) {
	exists, err := s.userRepo.ExistsByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	status := input.Status
	if status == 0 {
		status = models.UserStatusNormal
	}

	userType := input.UserType
	if userType == 0 {
		userType = models.UserTypeNormal
	}

	user := &models.User{
		Username:  input.Username,
		Password:  string(hashedPassword),
		RealName:  input.RealName,
		Phone:     input.Phone,
		Email:     input.Email,
		DeptID:    input.DeptID,
		UserType:  userType,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserInput 更新用户输入
type UpdateUserInput struct {
	RealName *string
	Phone    *string
	Email    *string
	DeptID   *int64
	Status   *int8
}

// Update 更新用户
func (s *UserService) Update(id int64, input *UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.RealName != nil {
		updates["real_name"] = *input.RealName
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.DeptID != nil {
		updates["dept_id"] = *input.DeptID
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if err := s.userRepo.Update(user, updates); err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(id)
}

// Delete 删除用户
func (s *UserService) Delete(id int64) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if user.Username == "admin" {
		return ErrCannotDeleteAdmin
	}

	return s.userRepo.Delete(id)
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(id int64, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.Update(user, map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	})
}
