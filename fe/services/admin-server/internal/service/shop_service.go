package service

import (
	"errors"
	"fmt"
	"time"

	"admin-server/internal/repository"
	"shared/models"
)

var (
	ErrShopNoExists = errors.New("店铺编号已存在")
)

// ShopService 店铺服务
type ShopService struct {
	shopRepo *repository.ShopRepository
}

// NewShopService 创建店铺服务
func NewShopService(shopRepo *repository.ShopRepository) *ShopService {
	return &ShopService{shopRepo: shopRepo}
}

// List 获取店铺列表
func (s *ShopService) List(shopNo, shopName string, merchantID *int64, status *int8, page, pageSize int) ([]models.Shop, int64, error) {
	offset := (page - 1) * pageSize
	return s.shopRepo.List(shopNo, shopName, merchantID, status, offset, pageSize)
}

// ListAll 获取所有店铺（下拉选择用）
func (s *ShopService) ListAll(merchantID *int64) ([]models.Shop, error) {
	status := int8(models.ShopStatusNormal)
	return s.shopRepo.ListAll(merchantID, &status)
}

// GetByID 获取店铺详情
func (s *ShopService) GetByID(id int64) (*models.Shop, error) {
	return s.shopRepo.FindByID(id)
}

// CreateShopInput 创建店铺输入
type CreateShopInput struct {
	ShopNo       string
	ShopName     string
	MerchantID   int64
	Building     *string
	Floor        *string
	RoomNo       *string
	ContactName  *string
	ContactPhone *string
	Status       int8
	Remark       *string
}

// Create 创建店铺
func (s *ShopService) Create(input *CreateShopInput) (*models.Shop, error) {
	// 生成店铺编号（如果未提供）
	shopNo := input.ShopNo
	if shopNo == "" {
		shopNo = fmt.Sprintf("S%s", time.Now().Format("20060102150405"))
	}

	// 检查店铺编号是否存在
	exists, err := s.shopRepo.ExistsByShopNo(shopNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrShopNoExists
	}

	status := input.Status
	if status == 0 {
		status = models.ShopStatusNormal
	}

	shop := &models.Shop{
		ShopNo:       shopNo,
		ShopName:     input.ShopName,
		MerchantID:   input.MerchantID,
		Building:     input.Building,
		Floor:        input.Floor,
		RoomNo:       input.RoomNo,
		ContactName:  input.ContactName,
		ContactPhone: input.ContactPhone,
		Status:       status,
		Remark:       input.Remark,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.shopRepo.Create(shop); err != nil {
		return nil, err
	}

	return shop, nil
}

// UpdateShopInput 更新店铺输入
type UpdateShopInput struct {
	ShopName     *string
	MerchantID   *int64
	Building     *string
	Floor        *string
	RoomNo       *string
	ContactName  *string
	ContactPhone *string
	Status       *int8
	Remark       *string
}

// Update 更新店铺
func (s *ShopService) Update(id int64, input *UpdateShopInput) (*models.Shop, error) {
	shop, err := s.shopRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.ShopName != nil {
		updates["shop_name"] = *input.ShopName
	}
	if input.MerchantID != nil {
		updates["merchant_id"] = *input.MerchantID
	}
	if input.Building != nil {
		updates["building"] = *input.Building
	}
	if input.Floor != nil {
		updates["floor"] = *input.Floor
	}
	if input.RoomNo != nil {
		updates["room_no"] = *input.RoomNo
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

	if err := s.shopRepo.Update(shop, updates); err != nil {
		return nil, err
	}

	return s.shopRepo.FindByID(id)
}

// Delete 删除店铺
func (s *ShopService) Delete(id int64) error {
	return s.shopRepo.Delete(id)
}

// ShopStats 店铺统计
type ShopStats struct {
	ElectricMeterCount int64 `json:"electric_meter_count"`
	WaterMeterCount    int64 `json:"water_meter_count"`
}

// GetStats 获取店铺统计信息
func (s *ShopService) GetStats(shopID int64) (*ShopStats, error) {
	electricMeterCount, waterMeterCount, err := s.shopRepo.GetStats(shopID)
	if err != nil {
		return nil, err
	}

	return &ShopStats{
		ElectricMeterCount: electricMeterCount,
		WaterMeterCount:    waterMeterCount,
	}, nil
}
