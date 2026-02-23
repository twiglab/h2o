package service

import (
	"errors"
	"time"

	"admin-server/internal/repository"
	"shared/models"

	"github.com/shopspring/decimal"
)

var (
	ErrRateCodeExists = errors.New("费率编码已存在")
	ErrRateNotFound   = errors.New("费率不存在")
	ErrRateInUse      = errors.New("该费率已被电表使用，无法删除")
)

// RateService 费率服务
type RateService struct {
	rateRepo  *repository.RateRepository
	meterRepo *repository.MeterRepository
}

// NewRateService 创建费率服务
func NewRateService(rateRepo *repository.RateRepository, meterRepo *repository.MeterRepository) *RateService {
	return &RateService{
		rateRepo:  rateRepo,
		meterRepo: meterRepo,
	}
}

// List 获取费率列表
func (s *RateService) List(keyword string, scope, status, calcMode *int8, page, pageSize int) ([]models.ElectricRate, int64, error) {
	offset := (page - 1) * pageSize
	return s.rateRepo.List(keyword, scope, status, calcMode, offset, pageSize)
}

// GetByID 获取费率详情
func (s *RateService) GetByID(id int64) (*models.ElectricRate, error) {
	return s.rateRepo.FindByID(id)
}

// TOUInput 分时电价输入
type TOUInput struct {
	PeriodName string
	StartTime  string
	EndTime    string
	UnitPrice  float64
}

// ServiceFeeInput 服务费输入
type ServiceFeeInput struct {
	FeeName  string
	FeeType  int8
	FeeValue float64
}

// CreateRateInput 创建费率输入
type CreateRateInput struct {
	RateCode      string
	RateName      string
	Scope         int8
	MerchantID    *int64
	CalcMode      int8
	UnitPrice     float64
	EffectiveDate *time.Time
	Status        int8
	Remark        *string
	TOUDetails    []TOUInput
	ServiceFees   []ServiceFeeInput
}

// Create 创建费率
func (s *RateService) Create(input *CreateRateInput) (*models.ElectricRate, error) {
	exists, err := s.rateRepo.ExistsByCode(input.RateCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrRateCodeExists
	}

	// 默认值处理
	calcMode := input.CalcMode
	if calcMode == 0 {
		calcMode = models.ElectricCalcModeFixed
	}
	scope := input.Scope
	if scope == 0 {
		scope = models.ElectricRateScopeGlobal
	}
	status := input.Status
	if status == 0 {
		status = models.ElectricRateStatusNormal
	}
	effectiveDate := time.Now()
	if input.EffectiveDate != nil {
		effectiveDate = *input.EffectiveDate
	}

	unitPrice := decimal.NewFromFloat(input.UnitPrice)

	rate := &models.ElectricRate{
		RateCode:      input.RateCode,
		RateName:      input.RateName,
		Scope:         scope,
		MerchantID:    input.MerchantID,
		CalcMode:      calcMode,
		UnitPrice:     &unitPrice,
		EffectiveDate: effectiveDate,
		Status:        status,
		Remark:        input.Remark,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tx := s.rateRepo.DB().Begin()

	if err := tx.Create(rate).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建分时电价详情
	for _, tou := range input.TOUDetails {
		touDetail := &models.ElectricRateTOU{
			RateID:     rate.ID,
			PeriodName: tou.PeriodName,
			StartTime:  tou.StartTime,
			EndTime:    tou.EndTime,
			UnitPrice:  decimal.NewFromFloat(tou.UnitPrice),
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(touDetail).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 创建服务费
	for _, sf := range input.ServiceFees {
		serviceFee := &models.ElectricRateServiceFee{
			RateID:    rate.ID,
			FeeName:   sf.FeeName,
			FeeType:   sf.FeeType,
			FeeValue:  decimal.NewFromFloat(sf.FeeValue),
			CreatedAt: time.Now(),
		}
		if err := tx.Create(serviceFee).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return rate, nil
}

// UpdateRateInput 更新费率输入
type UpdateRateInput struct {
	RateName  *string
	UnitPrice *float64
	Status    *int8
	Remark    *string
}

// Update 更新费率
func (s *RateService) Update(id int64, input *UpdateRateInput) (*models.ElectricRate, error) {
	rate, err := s.rateRepo.FindByID(id)
	if err != nil {
		return nil, ErrRateNotFound
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.RateName != nil {
		updates["rate_name"] = *input.RateName
	}
	if input.UnitPrice != nil {
		updates["unit_price"] = decimal.NewFromFloat(*input.UnitPrice)
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Remark != nil {
		updates["remark"] = *input.Remark
	}

	if err := s.rateRepo.Update(rate, updates); err != nil {
		return nil, err
	}

	return s.rateRepo.FindByID(id)
}

// Delete 删除费率
func (s *RateService) Delete(id int64) error {
	// 检查是否有电表使用此费率
	count, err := s.meterRepo.CountElectricByRateID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRateInUse
	}

	tx := s.rateRepo.DB().Begin()

	if err := s.rateRepo.DeleteTOUByRateID(id); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.rateRepo.DeleteServiceFeeByRateID(id); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.rateRepo.Delete(id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
