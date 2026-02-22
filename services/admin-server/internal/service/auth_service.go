package service

import (
	"errors"

	"admin-server/internal/repository"
	"shared/middleware"
	"shared/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrOldPasswordWrong   = errors.New("原密码错误")
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	jwt      *middleware.JWTMiddleware
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, jwt *middleware.JWTMiddleware) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *models.User
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindActiveByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.jwt.GenerateTokenSimple(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := s.jwt.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return s.jwt.GenerateTokenSimple(claims.UserID, claims.Username)
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(userID int64) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrOldPasswordWrong
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.Update(user, map[string]interface{}{
		"password": string(hashedPassword),
	})
}

// JWT 返回JWT中间件
func (s *AuthService) JWT() *middleware.JWTMiddleware {
	return s.jwt
}
