package service

import (
	"errors"
	"fmt"
	"time"

	"admin-server/internal/repository"
	"shared/models"
)

var (
	ErrMerchantNoExists = errors.New("商户编号已存在")
)

// MerchantService 商户服务
type MerchantService struct {
	merchantRepo *repository.MerchantRepository
}

// NewMerchantService 创建商户服务
func NewMerchantService(merchantRepo *repository.MerchantRepository) *MerchantService {
	return &MerchantService{merchantRepo: merchantRepo}
}

// List 获取商户列表
func (s *MerchantService) List(merchantNo, merchantName string, status *int8, page, pageSize int) ([]models.Merchant, int64, error) {
	offset := (page - 1) * pageSize
	return s.merchantRepo.List(merchantNo, merchantName, status, offset, pageSize)
}

// ListAll 获取所有商户（下拉选择用）
func (s *MerchantService) ListAll() ([]models.Merchant, error) {
	status := int8(models.MerchantStatusNormal)
	return s.merchantRepo.ListAll(&status)
}

// GetByID 获取商户详情
func (s *MerchantService) GetByID(id int64) (*models.Merchant, error) {
	return s.merchantRepo.FindByID(id)
}

// CreateMerchantInput 创建商户输入
type CreateMerchantInput struct {
	MerchantName string
	MerchantType int8
	ContactName  *string
	ContactPhone *string
	Status       int8
	Remark       *string
}

// Create 创建商户
func (s *MerchantService) Create(input *CreateMerchantInput) (*models.Merchant, error) {
	// 自动生成商户编号
	merchantNo := fmt.Sprintf("M%s", time.Now().Format("20060102150405"))

	// 检查商户编号是否存在
	exists, err := s.merchantRepo.ExistsByMerchantNo(merchantNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrMerchantNoExists
	}

	status := input.Status
	if status == 0 {
		status = models.MerchantStatusNormal
	}

	merchantType := input.MerchantType
	if merchantType == 0 {
		merchantType = models.MerchantTypeEnterprise
	}

	merchant := &models.Merchant{
		MerchantNo:   merchantNo,
		MerchantName: input.MerchantName,
		MerchantType: merchantType,
		ContactName:  input.ContactName,
		ContactPhone: input.ContactPhone,
		Status:       status,
		Remark:       input.Remark,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.merchantRepo.Create(merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

// UpdateMerchantInput 更新商户输入
type UpdateMerchantInput struct {
	MerchantName *string
	MerchantType *int8
	ContactName  *string
	ContactPhone *string
	Status       *int8
	Remark       *string
}

// Update 更新商户
func (s *MerchantService) Update(id int64, input *UpdateMerchantInput) (*models.Merchant, error) {
	merchant, err := s.merchantRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.MerchantName != nil {
		updates["merchant_name"] = *input.MerchantName
	}
	if input.MerchantType != nil {
		updates["merchant_type"] = *input.MerchantType
	}
	if input.ContactName != nil {
		updates["contact_name"] = *input.ContactName
	}
	if input.ContactPhone != nil {
		updates["contact_phone"] = *input.ContactPhone
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Remark != nil {
		updates["remark"] = *input.Remark
	}

	if err := s.merchantRepo.Update(merchant, updates); err != nil {
		return nil, err
	}

	return s.merchantRepo.FindByID(id)
}

// Delete 删除商户
func (s *MerchantService) Delete(id int64) error {
	return s.merchantRepo.Delete(id)
}

// MerchantStats 商户统计
type MerchantStats struct {
	ShopCount          int64 `json:"shop_count"`
	ElectricMeterCount int64 `json:"electric_meter_count"`
	WaterMeterCount    int64 `json:"water_meter_count"`
}

// GetStats 获取商户统计信息
func (s *MerchantService) GetStats(merchantID int64) (*MerchantStats, error) {
	shopCount, electricMeterCount, waterMeterCount, err := s.merchantRepo.GetStats(merchantID)
	if err != nil {
		return nil, err
	}

	return &MerchantStats{
		ShopCount:          shopCount,
		ElectricMeterCount: electricMeterCount,
		WaterMeterCount:    waterMeterCount,
	}, nil
}
