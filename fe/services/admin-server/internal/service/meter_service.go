package service

import (
	"errors"
	"time"

	"admin-server/internal/repository"
	"shared/models"
	"shared/utils"

	"github.com/shopspring/decimal"
)

var (
	ErrMeterNoExists = errors.New("电表编号已存在")
	ErrMeterNotFound = errors.New("电表不存在")
	ErrReadingTooLow = errors.New("读数不能小于当前读数")
)

// ReadingInfo 读数信息（统一响应结构）
type ReadingInfo struct {
	ID           int64           `json:"id"`
	MeterID      int64           `json:"meter_id"`
	MeterNo      string          `json:"meter_no"`
	ReadingValue decimal.Decimal `json:"reading_value"`
	ReadingTime  time.Time       `json:"reading_time"`
	CollectType  int8            `json:"collect_type"`
	Status       int8            `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
}

// MeterInfo 表计信息（统一响应结构）
type MeterInfo struct {
	ID             int64           `json:"id"`
	MeterNo        string          `json:"meter_no"`
	MeterType      int16           `json:"meter_type"`
	MerchantID     int64           `json:"merchant_id"`
	MerchantName   string          `json:"merchant_name,omitempty"`
	ShopID         *int64          `json:"shop_id"`
	ShopName       string          `json:"shop_name,omitempty"`
	RateID         *int64          `json:"rate_id"`
	RateName       string          `json:"rate_name,omitempty"`
	CommAddr       *string         `json:"comm_addr"`
	Multiplier     decimal.Decimal `json:"multiplier"`
	InitReading    decimal.Decimal `json:"init_reading"`
	CurrentReading decimal.Decimal `json:"current_reading"`
	LastCollectAt  *time.Time      `json:"last_collect_at"`
	OnlineStatus   int8            `json:"online_status"`
	Status         int8            `json:"status"`
	InstallDate    *time.Time      `json:"install_date"`
	Remark         *string         `json:"remark"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// MeterService 表计服务
type MeterService struct {
	meterRepo   *repository.MeterRepository
	readingRepo *repository.ReadingRepository
}

// NewMeterService 创建表计服务
func NewMeterService(meterRepo *repository.MeterRepository, readingRepo *repository.ReadingRepository) *MeterService {
	return &MeterService{
		meterRepo:   meterRepo,
		readingRepo: readingRepo,
	}
}

// ListElectric 获取电表列表
func (s *MeterService) ListElectric(meterNo string, status, onlineStatus *int16, page, pageSize int) ([]MeterInfo, int64, error) {
	offset := (page - 1) * pageSize
	meters, total, err := s.meterRepo.ListElectric(meterNo, status, onlineStatus, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []MeterInfo
	for _, m := range meters {
		info := MeterInfo{
			ID:             m.ID,
			MeterNo:        m.MeterNo,
			MeterType:      models.MeterTypeElectric,
			MerchantID:     m.MerchantID,
			ShopID:         m.ShopID,
			RateID:         m.RateID,
			CommAddr:       m.CommAddr,
			Multiplier:     m.Multiplier,
			InitReading:    m.InitReading,
			CurrentReading: m.CurrentReading,
			LastCollectAt:  m.LastCollectAt,
			OnlineStatus:   m.OnlineStatus,
			Status:         m.Status,
			InstallDate:    m.InstallDate,
			Remark:         m.Remark,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		}
		if m.Merchant != nil {
			info.MerchantName = m.Merchant.MerchantName
		}
		if m.Shop != nil {
			info.ShopName = m.Shop.ShopName
		}
		if m.ElectricRate != nil {
			info.RateName = m.ElectricRate.RateName
		}
		result = append(result, info)
	}

	return result, total, nil
}

// ListWater 获取水表列表
func (s *MeterService) ListWater(meterNo string, status, onlineStatus *int16, page, pageSize int) ([]MeterInfo, int64, error) {
	offset := (page - 1) * pageSize
	meters, total, err := s.meterRepo.ListWater(meterNo, status, onlineStatus, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []MeterInfo
	for _, m := range meters {
		info := MeterInfo{
			ID:             m.ID,
			MeterNo:        m.MeterNo,
			MeterType:      models.MeterTypeWater,
			MerchantID:     m.MerchantID,
			ShopID:         m.ShopID,
			RateID:         m.RateID,
			CommAddr:       m.CommAddr,
			Multiplier:     m.Multiplier,
			InitReading:    m.InitReading,
			CurrentReading: m.CurrentReading,
			LastCollectAt:  m.LastCollectAt,
			OnlineStatus:   m.OnlineStatus,
			Status:         m.Status,
			InstallDate:    m.InstallDate,
			Remark:         m.Remark,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		}
		if m.Merchant != nil {
			info.MerchantName = m.Merchant.MerchantName
		}
		if m.Shop != nil {
			info.ShopName = m.Shop.ShopName
		}
		if m.WaterRate != nil {
			info.RateName = m.WaterRate.RateName
		}
		result = append(result, info)
	}

	return result, total, nil
}

// GetElectricByID 获取电表详情
func (s *MeterService) GetElectricByID(id int64) (*MeterInfo, []ReadingInfo, error) {
	meter, err := s.meterRepo.FindElectricByID(id)
	if err != nil {
		return nil, nil, err
	}

	info := &MeterInfo{
		ID:             meter.ID,
		MeterNo:        meter.MeterNo,
		MeterType:      models.MeterTypeElectric,
		MerchantID:     meter.MerchantID,
		ShopID:         meter.ShopID,
		RateID:         meter.RateID,
		CommAddr:       meter.CommAddr,
		Multiplier:     meter.Multiplier,
		InitReading:    meter.InitReading,
		CurrentReading: meter.CurrentReading,
		LastCollectAt:  meter.LastCollectAt,
		OnlineStatus:   meter.OnlineStatus,
		Status:         meter.Status,
		InstallDate:    meter.InstallDate,
		Remark:         meter.Remark,
		CreatedAt:      meter.CreatedAt,
		UpdatedAt:      meter.UpdatedAt,
	}
	if meter.Merchant != nil {
		info.MerchantName = meter.Merchant.MerchantName
	}
	if meter.Shop != nil {
		info.ShopName = meter.Shop.ShopName
	}
	if meter.ElectricRate != nil {
		info.RateName = meter.ElectricRate.RateName
	}

	electricReadings, err := s.readingRepo.ListRecentElectricByMeterID(id, 10)
	if err != nil {
		return nil, nil, err
	}

	var readings []ReadingInfo
	for _, r := range electricReadings {
		readings = append(readings, ReadingInfo{
			ID:           r.ID,
			MeterID:      r.MeterID,
			MeterNo:      r.MeterNo,
			ReadingValue: r.ReadingValue,
			ReadingTime:  r.ReadingTime,
			CollectType:  r.CollectType,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		})
	}

	return info, readings, nil
}

// GetWaterByID 获取水表详情
func (s *MeterService) GetWaterByID(id int64) (*MeterInfo, []ReadingInfo, error) {
	meter, err := s.meterRepo.FindWaterByID(id)
	if err != nil {
		return nil, nil, err
	}

	info := &MeterInfo{
		ID:             meter.ID,
		MeterNo:        meter.MeterNo,
		MeterType:      models.MeterTypeWater,
		MerchantID:     meter.MerchantID,
		ShopID:         meter.ShopID,
		RateID:         meter.RateID,
		CommAddr:       meter.CommAddr,
		Multiplier:     meter.Multiplier,
		InitReading:    meter.InitReading,
		CurrentReading: meter.CurrentReading,
		LastCollectAt:  meter.LastCollectAt,
		OnlineStatus:   meter.OnlineStatus,
		Status:         meter.Status,
		InstallDate:    meter.InstallDate,
		Remark:         meter.Remark,
		CreatedAt:      meter.CreatedAt,
		UpdatedAt:      meter.UpdatedAt,
	}
	if meter.Merchant != nil {
		info.MerchantName = meter.Merchant.MerchantName
	}
	if meter.Shop != nil {
		info.ShopName = meter.Shop.ShopName
	}
	if meter.WaterRate != nil {
		info.RateName = meter.WaterRate.RateName
	}

	waterReadings, err := s.readingRepo.ListRecentWaterByMeterID(id, 10)
	if err != nil {
		return nil, nil, err
	}

	var readings []ReadingInfo
	for _, r := range waterReadings {
		readings = append(readings, ReadingInfo{
			ID:           r.ID,
			MeterID:      r.MeterID,
			MeterNo:      r.MeterNo,
			ReadingValue: r.ReadingValue,
			ReadingTime:  r.ReadingTime,
			CollectType:  r.CollectType,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
		})
	}

	return info, readings, nil
}

// CreateMeterInput 创建表计输入
type CreateMeterInput struct {
	MeterNo     string
	MeterType   int16
	RateID      *int64
	CommAddr    *string
	Multiplier  float64
	InitReading float64
	Remark      *string
}

// CreateElectric 创建电表
func (s *MeterService) CreateElectric(input *CreateMeterInput) (*MeterInfo, error) {
	exists, err := s.meterRepo.ExistsElectricByMeterNo(input.MeterNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrMeterNoExists
	}

	multiplier := decimal.NewFromFloat(1.0)
	if input.Multiplier > 0 {
		multiplier = decimal.NewFromFloat(input.Multiplier)
	}

	meter := &models.ElectricMeter{
		MeterNo:        input.MeterNo,
		RateID:         input.RateID,
		CommAddr:       input.CommAddr,
		Multiplier:     multiplier,
		InitReading:    decimal.NewFromFloat(input.InitReading),
		CurrentReading: decimal.NewFromFloat(input.InitReading),
		OnlineStatus:   models.MeterOffline,
		Status:         models.MeterStatusNormal,
		Remark:         input.Remark,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.meterRepo.CreateElectric(meter); err != nil {
		return nil, err
	}

	return &MeterInfo{
		ID:             meter.ID,
		MeterNo:        meter.MeterNo,
		MeterType:      models.MeterTypeElectric,
		RateID:         meter.RateID,
		CommAddr:       meter.CommAddr,
		Multiplier:     meter.Multiplier,
		InitReading:    meter.InitReading,
		CurrentReading: meter.CurrentReading,
		OnlineStatus:   meter.OnlineStatus,
		Status:         meter.Status,
		CreatedAt:      meter.CreatedAt,
		UpdatedAt:      meter.UpdatedAt,
	}, nil
}

// CreateWater 创建水表
func (s *MeterService) CreateWater(input *CreateMeterInput) (*MeterInfo, error) {
	exists, err := s.meterRepo.ExistsWaterByMeterNo(input.MeterNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrMeterNoExists
	}

	multiplier := decimal.NewFromFloat(1.0)
	if input.Multiplier > 0 {
		multiplier = decimal.NewFromFloat(input.Multiplier)
	}

	meter := &models.WaterMeter{
		MeterNo:        input.MeterNo,
		RateID:         input.RateID,
		CommAddr:       input.CommAddr,
		Multiplier:     multiplier,
		InitReading:    decimal.NewFromFloat(input.InitReading),
		CurrentReading: decimal.NewFromFloat(input.InitReading),
		OnlineStatus:   models.MeterOffline,
		Status:         models.MeterStatusNormal,
		Remark:         input.Remark,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.meterRepo.CreateWater(meter); err != nil {
		return nil, err
	}

	return &MeterInfo{
		ID:             meter.ID,
		MeterNo:        meter.MeterNo,
		MeterType:      models.MeterTypeWater,
		RateID:         meter.RateID,
		CommAddr:       meter.CommAddr,
		Multiplier:     meter.Multiplier,
		InitReading:    meter.InitReading,
		CurrentReading: meter.CurrentReading,
		OnlineStatus:   meter.OnlineStatus,
		Status:         meter.Status,
		CreatedAt:      meter.CreatedAt,
		UpdatedAt:      meter.UpdatedAt,
	}, nil
}

// UpdateMeterInput 更新表计输入
type UpdateMeterInput struct {
	MeterNo    *string
	RateID     *int64
	CommAddr   *string
	Multiplier *float64
	Status     *int16
	Remark     *string
}

// UpdateElectric 更新电表
func (s *MeterService) UpdateElectric(id int64, input *UpdateMeterInput) (*MeterInfo, error) {
	meter, err := s.meterRepo.FindElectricByID(id)
	if err != nil {
		return nil, ErrMeterNotFound
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.MeterNo != nil {
		updates["meter_no"] = *input.MeterNo
	}
	if input.RateID != nil {
		updates["rate_id"] = *input.RateID
	}
	if input.CommAddr != nil {
		updates["comm_addr"] = *input.CommAddr
	}
	if input.Multiplier != nil {
		updates["multiplier"] = decimal.NewFromFloat(*input.Multiplier)
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Remark != nil {
		updates["remark"] = *input.Remark
	}

	if err := s.meterRepo.UpdateElectric(meter, updates); err != nil {
		return nil, err
	}

	return s.GetElectricByIDSimple(id)
}

// UpdateWater 更新水表
func (s *MeterService) UpdateWater(id int64, input *UpdateMeterInput) (*MeterInfo, error) {
	meter, err := s.meterRepo.FindWaterByID(id)
	if err != nil {
		return nil, ErrMeterNotFound
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.MeterNo != nil {
		updates["meter_no"] = *input.MeterNo
	}
	if input.RateID != nil {
		updates["rate_id"] = *input.RateID
	}
	if input.CommAddr != nil {
		updates["comm_addr"] = *input.CommAddr
	}
	if input.Multiplier != nil {
		updates["multiplier"] = decimal.NewFromFloat(*input.Multiplier)
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Remark != nil {
		updates["remark"] = *input.Remark
	}

	if err := s.meterRepo.UpdateWater(meter, updates); err != nil {
		return nil, err
	}

	return s.GetWaterByIDSimple(id)
}

// GetElectricByIDSimple 获取电表基本信息
func (s *MeterService) GetElectricByIDSimple(id int64) (*MeterInfo, error) {
	meter, err := s.meterRepo.FindElectricByID(id)
	if err != nil {
		return nil, err
	}

	info := &MeterInfo{
		ID:             meter.ID,
		MeterNo:        meter.MeterNo,
		MeterType:      models.MeterTypeElectric,
		MerchantID:     meter.MerchantID,
		ShopID:         meter.ShopID,
		RateID:         meter.RateID,
		CommAddr:       meter.CommAddr,
		Multiplier:     meter.Multiplier,
		InitReading:    meter.InitReading,
		CurrentReading: meter.CurrentReading,
		LastCollectAt:  meter.LastCollectAt,
		OnlineStatus:   meter.OnlineStatus,
		Status:         meter.Status,
		InstallDate:    meter.InstallDate,
		Remark:         meter.Remark,
		CreatedAt:      meter.CreatedAt,
		UpdatedAt:      meter.UpdatedAt,
	}
	if meter.Merchant != nil {
		info.MerchantName = meter.Merchant.MerchantName
	}
	if meter.Shop != nil {
		info.ShopName = meter.Shop.ShopName
	}
	if meter.ElectricRate != nil {
		info.RateName = meter.ElectricRate.RateName
	}

	return info, nil
}

// GetWaterByIDSimple 获取水表基本信息
func (s *MeterService) GetWaterByIDSimple(id int64) (*MeterInfo, error) {
	meter, err := s.meterRepo.FindWaterByID(id)
	if err != nil {
		return nil, err
	}

	info := &MeterInfo{
		ID:             meter.ID,
		MeterNo:        meter.MeterNo,
		MeterType:      models.MeterTypeWater,
		MerchantID:     meter.MerchantID,
		ShopID:         meter.ShopID,
		RateID:         meter.RateID,
		CommAddr:       meter.CommAddr,
		Multiplier:     meter.Multiplier,
		InitReading:    meter.InitReading,
		CurrentReading: meter.CurrentReading,
		LastCollectAt:  meter.LastCollectAt,
		OnlineStatus:   meter.OnlineStatus,
		Status:         meter.Status,
		InstallDate:    meter.InstallDate,
		Remark:         meter.Remark,
		CreatedAt:      meter.CreatedAt,
		UpdatedAt:      meter.UpdatedAt,
	}
	if meter.Merchant != nil {
		info.MerchantName = meter.Merchant.MerchantName
	}
	if meter.Shop != nil {
		info.ShopName = meter.Shop.ShopName
	}
	if meter.WaterRate != nil {
		info.RateName = meter.WaterRate.RateName
	}

	return info, nil
}

// DeleteElectric 删除电表
func (s *MeterService) DeleteElectric(id int64) error {
	return s.meterRepo.DeleteElectric(id)
}

// DeleteWater 删除水表
func (s *MeterService) DeleteWater(id int64) error {
	return s.meterRepo.DeleteWater(id)
}

// ManualReadingInput 手动抄表输入
type ManualReadingInput struct {
	MeterID      int64
	MeterType    int16
	ReadingValue float64
	ReadingTime  *time.Time
	OperatorID   *int64
	OperatorName *string
	Remark       *string
}

// ManualReading 手动抄表
func (s *MeterService) ManualReading(input *ManualReadingInput) error {
	readingValue := decimal.NewFromFloat(input.ReadingValue)
	readingTime := time.Now()
	if input.ReadingTime != nil {
		readingTime = *input.ReadingTime
	}

	tx := s.readingRepo.DB().Begin()

	if input.MeterType == models.MeterTypeElectric {
		meter, err := s.meterRepo.FindElectricByID(input.MeterID)
		if err != nil {
			return ErrMeterNotFound
		}

		if readingValue.LessThan(meter.CurrentReading) {
			return ErrReadingTooLow
		}

		// 创建电表手工抄表记录
		manualReading := &models.ElectricManualReading{
			ReadingNo:       utils.GenerateNo("MR"),
			MerchantID:      meter.MerchantID,
			MeterID:         meter.ID,
			MeterNo:         meter.MeterNo,
			PreviousReading: meter.CurrentReading,
			CurrentReading:  readingValue,
			ReadingTime:     readingTime,
			OperatorID:      input.OperatorID,
			OperatorName:    input.OperatorName,
			Status:          1,
			Remark:          input.Remark,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := s.readingRepo.CreateElectricManualReading(tx, manualReading); err != nil {
			tx.Rollback()
			return err
		}

		// 创建电表读数记录
		electricReading := &models.ElectricReading{
			MerchantID:   meter.MerchantID,
			MeterID:      meter.ID,
			MeterNo:      meter.MeterNo,
			ReadingValue: readingValue,
			ReadingTime:  readingTime,
			CollectType:  models.CollectTypeManual,
			OperatorID:   input.OperatorID,
			Status:       models.ReadingStatusNormal,
			CreatedAt:    time.Now(),
		}
		if err := s.readingRepo.CreateElectricReading(tx, electricReading); err != nil {
			tx.Rollback()
			return err
		}

		// 更新电表当前读数
		if err := s.meterRepo.UpdateElectricReading(tx, meter.ID, readingValue, readingTime); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		meter, err := s.meterRepo.FindWaterByID(input.MeterID)
		if err != nil {
			return ErrMeterNotFound
		}

		if readingValue.LessThan(meter.CurrentReading) {
			return ErrReadingTooLow
		}

		// 创建水表手工抄表记录
		manualReading := &models.WaterManualReading{
			ReadingNo:       utils.GenerateNo("MR"),
			MerchantID:      meter.MerchantID,
			MeterID:         meter.ID,
			MeterNo:         meter.MeterNo,
			PreviousReading: meter.CurrentReading,
			CurrentReading:  readingValue,
			ReadingTime:     readingTime,
			OperatorID:      input.OperatorID,
			OperatorName:    input.OperatorName,
			Status:          1,
			Remark:          input.Remark,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := s.readingRepo.CreateWaterManualReading(tx, manualReading); err != nil {
			tx.Rollback()
			return err
		}

		// 创建水表读数记录
		waterReading := &models.WaterReading{
			MerchantID:   meter.MerchantID,
			MeterID:      meter.ID,
			MeterNo:      meter.MeterNo,
			ReadingValue: readingValue,
			ReadingTime:  readingTime,
			CollectType:  models.CollectTypeManual,
			OperatorID:   input.OperatorID,
			Status:       models.ReadingStatusNormal,
			CreatedAt:    time.Now(),
		}
		if err := s.readingRepo.CreateWaterReading(tx, waterReading); err != nil {
			tx.Rollback()
			return err
		}

		// 更新水表当前读数
		if err := s.meterRepo.UpdateWaterReading(tx, meter.ID, readingValue, readingTime); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetReadings 获取表计读数历史
func (s *MeterService) GetReadings(meterID int64, meterType int16, startDate, endDate string, page, pageSize int) ([]ReadingInfo, int64, error) {
	offset := (page - 1) * pageSize
	var readings []ReadingInfo
	var total int64

	if meterType == models.MeterTypeElectric {
		electricReadings, count, err := s.readingRepo.ListElectricByMeterID(meterID, startDate, endDate, offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
		total = count
		for _, r := range electricReadings {
			readings = append(readings, ReadingInfo{
				ID:           r.ID,
				MeterID:      r.MeterID,
				MeterNo:      r.MeterNo,
				ReadingValue: r.ReadingValue,
				ReadingTime:  r.ReadingTime,
				CollectType:  r.CollectType,
				Status:       r.Status,
				CreatedAt:    r.CreatedAt,
			})
		}
	} else {
		waterReadings, count, err := s.readingRepo.ListWaterByMeterID(meterID, startDate, endDate, offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
		total = count
		for _, r := range waterReadings {
			readings = append(readings, ReadingInfo{
				ID:           r.ID,
				MeterID:      r.MeterID,
				MeterNo:      r.MeterNo,
				ReadingValue: r.ReadingValue,
				ReadingTime:  r.ReadingTime,
				CollectType:  r.CollectType,
				Status:       r.Status,
				CreatedAt:    r.CreatedAt,
			})
		}
	}

	return readings, total, nil
}

// ManualReadingInfo 手工抄表记录信息（统一响应结构）
type ManualReadingInfo struct {
	ID              int64           `json:"id"`
	ReadingNo       string          `json:"reading_no"`
	MeterID         int64           `json:"meter_id"`
	MeterNo         string          `json:"meter_no"`
	MeterType       int16           `json:"meter_type"`
	PreviousReading decimal.Decimal `json:"previous_reading"`
	CurrentReading  decimal.Decimal `json:"current_reading"`
	ReadingTime     time.Time       `json:"reading_time"`
	OperatorID      *int64          `json:"operator_id"`
	OperatorName    *string         `json:"operator_name"`
	Status          int8            `json:"status"`
	Remark          *string         `json:"remark"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ListManualReadings 获取手工抄表记录列表
func (s *MeterService) ListManualReadings(meterType *int16, meterNo, startDate, endDate string, page, pageSize int) ([]ManualReadingInfo, int64, error) {
	offset := (page - 1) * pageSize
	var result []ManualReadingInfo
	var total int64

	// 如果指定了表类型，只查询对应类型
	if meterType != nil {
		if *meterType == models.MeterTypeElectric {
			readings, count, err := s.readingRepo.ListElectricManualReadings(meterNo, startDate, endDate, offset, pageSize)
			if err != nil {
				return nil, 0, err
			}
			total = count
			for _, r := range readings {
				result = append(result, ManualReadingInfo{
					ID:              r.ID,
					ReadingNo:       r.ReadingNo,
					MeterID:         r.MeterID,
					MeterNo:         r.MeterNo,
					MeterType:       models.MeterTypeElectric,
					PreviousReading: r.PreviousReading,
					CurrentReading:  r.CurrentReading,
					ReadingTime:     r.ReadingTime,
					OperatorID:      r.OperatorID,
					OperatorName:    r.OperatorName,
					Status:          r.Status,
					Remark:          r.Remark,
					CreatedAt:       r.CreatedAt,
				})
			}
		} else {
			readings, count, err := s.readingRepo.ListWaterManualReadings(meterNo, startDate, endDate, offset, pageSize)
			if err != nil {
				return nil, 0, err
			}
			total = count
			for _, r := range readings {
				result = append(result, ManualReadingInfo{
					ID:              r.ID,
					ReadingNo:       r.ReadingNo,
					MeterID:         r.MeterID,
					MeterNo:         r.MeterNo,
					MeterType:       models.MeterTypeWater,
					PreviousReading: r.PreviousReading,
					CurrentReading:  r.CurrentReading,
					ReadingTime:     r.ReadingTime,
					OperatorID:      r.OperatorID,
					OperatorName:    r.OperatorName,
					Status:          r.Status,
					Remark:          r.Remark,
					CreatedAt:       r.CreatedAt,
				})
			}
		}
	} else {
		// 没有指定类型，查询所有（先查电表，再查水表）
		electricReadings, electricCount, err := s.readingRepo.ListElectricManualReadings(meterNo, startDate, endDate, offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
		for _, r := range electricReadings {
			result = append(result, ManualReadingInfo{
				ID:              r.ID,
				ReadingNo:       r.ReadingNo,
				MeterID:         r.MeterID,
				MeterNo:         r.MeterNo,
				MeterType:       models.MeterTypeElectric,
				PreviousReading: r.PreviousReading,
				CurrentReading:  r.CurrentReading,
				ReadingTime:     r.ReadingTime,
				OperatorID:      r.OperatorID,
				OperatorName:    r.OperatorName,
				Status:          r.Status,
				Remark:          r.Remark,
				CreatedAt:       r.CreatedAt,
			})
		}

		waterReadings, waterCount, err := s.readingRepo.ListWaterManualReadings(meterNo, startDate, endDate, offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
		for _, r := range waterReadings {
			result = append(result, ManualReadingInfo{
				ID:              r.ID,
				ReadingNo:       r.ReadingNo,
				MeterID:         r.MeterID,
				MeterNo:         r.MeterNo,
				MeterType:       models.MeterTypeWater,
				PreviousReading: r.PreviousReading,
				CurrentReading:  r.CurrentReading,
				ReadingTime:     r.ReadingTime,
				OperatorID:      r.OperatorID,
				OperatorName:    r.OperatorName,
				Status:          r.Status,
				Remark:          r.Remark,
				CreatedAt:       r.CreatedAt,
			})
		}

		total = electricCount + waterCount
	}

	return result, total, nil
}
