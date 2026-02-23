package handler

import (
	"time"

	"admin-server/internal/service"
	"shared/middleware"
	"shared/models"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type MeterHandler struct {
	meterService *service.MeterService
}

func NewMeterHandler(meterService *service.MeterService) *MeterHandler {
	return &MeterHandler{meterService: meterService}
}

// ListElectric 获取电表列表
func (h *MeterHandler) ListElectric(c *gin.Context) {
	page, pageSize := utils.GetPage(c)
	meterNo := c.Query("meter_no")
	status := utils.ParseInt16Query(c, "status")
	onlineStatus := utils.ParseInt16Query(c, "online_status")

	meters, total, err := h.meterService.ListElectric(meterNo, status, onlineStatus, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, meters, total, page, pageSize)
}

// ListWater 获取水表列表
func (h *MeterHandler) ListWater(c *gin.Context) {
	page, pageSize := utils.GetPage(c)
	meterNo := c.Query("meter_no")
	status := utils.ParseInt16Query(c, "status")
	onlineStatus := utils.ParseInt16Query(c, "online_status")

	meters, total, err := h.meterService.ListWater(meterNo, status, onlineStatus, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, meters, total, page, pageSize)
}

// GetElectric 获取电表详情
func (h *MeterHandler) GetElectric(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	meter, readings, err := h.meterService.GetElectricByID(id)
	if err != nil {
		response.NotFound(c, "电表不存在")
		return
	}

	response.Success(c, map[string]interface{}{
		"meter":    meter,
		"readings": readings,
	})
}

// GetWater 获取水表详情
func (h *MeterHandler) GetWater(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	meter, readings, err := h.meterService.GetWaterByID(id)
	if err != nil {
		response.NotFound(c, "水表不存在")
		return
	}

	response.Success(c, map[string]interface{}{
		"meter":    meter,
		"readings": readings,
	})
}

type CreateMeterRequest struct {
	MeterNo     string  `json:"meter_no" binding:"required"`
	RateID      *int64  `json:"rate_id"`
	CommAddr    *string `json:"comm_addr"`
	Multiplier  float64 `json:"multiplier"`
	InitReading float64 `json:"init_reading"`
	Remark      *string `json:"remark"`
}

// CreateElectric 创建电表
func (h *MeterHandler) CreateElectric(c *gin.Context) {
	var req CreateMeterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.CreateMeterInput{
		MeterNo:     req.MeterNo,
		MeterType:   models.MeterTypeElectric,
		RateID:      req.RateID,
		CommAddr:    req.CommAddr,
		Multiplier:  req.Multiplier,
		InitReading: req.InitReading,
		Remark:      req.Remark,
	}

	meter, err := h.meterService.CreateElectric(input)
	if err != nil {
		if err == service.ErrMeterNoExists {
			response.Error(c, 400, "电表编号已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, meter)
}

// CreateWater 创建水表
func (h *MeterHandler) CreateWater(c *gin.Context) {
	var req CreateMeterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.CreateMeterInput{
		MeterNo:     req.MeterNo,
		MeterType:   models.MeterTypeWater,
		RateID:      req.RateID,
		CommAddr:    req.CommAddr,
		Multiplier:  req.Multiplier,
		InitReading: req.InitReading,
		Remark:      req.Remark,
	}

	meter, err := h.meterService.CreateWater(input)
	if err != nil {
		if err == service.ErrMeterNoExists {
			response.Error(c, 400, "水表编号已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, meter)
}

type UpdateMeterRequest struct {
	MeterNo    *string  `json:"meter_no"`
	RateID     *int64   `json:"rate_id"`
	CommAddr   *string  `json:"comm_addr"`
	Multiplier *float64 `json:"multiplier"`
	Status     *int16   `json:"status"`
	Remark     *string  `json:"remark"`
}

// UpdateElectric 更新电表
func (h *MeterHandler) UpdateElectric(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateMeterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateMeterInput{
		MeterNo:    req.MeterNo,
		RateID:     req.RateID,
		CommAddr:   req.CommAddr,
		Multiplier: req.Multiplier,
		Status:     req.Status,
		Remark:     req.Remark,
	}

	meter, err := h.meterService.UpdateElectric(id, input)
	if err != nil {
		if err == service.ErrMeterNotFound {
			response.NotFound(c, "电表不存在")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, meter)
}

// UpdateWater 更新水表
func (h *MeterHandler) UpdateWater(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateMeterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateMeterInput{
		MeterNo:    req.MeterNo,
		RateID:     req.RateID,
		CommAddr:   req.CommAddr,
		Multiplier: req.Multiplier,
		Status:     req.Status,
		Remark:     req.Remark,
	}

	meter, err := h.meterService.UpdateWater(id, input)
	if err != nil {
		if err == service.ErrMeterNotFound {
			response.NotFound(c, "水表不存在")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, meter)
}

// DeleteElectric 删除电表
func (h *MeterHandler) DeleteElectric(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.meterService.DeleteElectric(id); err != nil {
		response.ServerError(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// DeleteWater 删除水表
func (h *MeterHandler) DeleteWater(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.meterService.DeleteWater(id); err != nil {
		response.ServerError(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

type ManualReadingRequest struct {
	MeterType    int16   `json:"meter_type" binding:"required"`
	ReadingValue float64 `json:"reading_value" binding:"required"`
	ReadingTime  string  `json:"reading_time"`
	Remark       *string `json:"remark"`
}

// ManualReadingElectric 电表手动抄表
func (h *MeterHandler) ManualReadingElectric(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req ManualReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var readingTime *time.Time
	if req.ReadingTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.ReadingTime); err == nil {
			readingTime = &t
		}
	}

	// 获取操作员信息
	var operatorID *int64
	var operatorName *string
	if userID, exists := c.Get("userID"); exists {
		uid := userID.(int64)
		operatorID = &uid
		if username := middleware.GetUsername(c); username != "" {
			operatorName = &username
		}
	}

	input := &service.ManualReadingInput{
		MeterID:      id,
		MeterType:    models.MeterTypeElectric,
		ReadingValue: req.ReadingValue,
		ReadingTime:  readingTime,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Remark:       req.Remark,
	}

	if err := h.meterService.ManualReading(input); err != nil {
		if err == service.ErrMeterNotFound {
			response.NotFound(c, "电表不存在")
			return
		}
		if err == service.ErrReadingTooLow {
			response.Error(c, 400, "读数不能小于当前读数")
			return
		}
		response.ServerError(c, "抄表失败")
		return
	}

	response.SuccessWithMessage(c, "抄表成功", nil)
}

// ManualReadingWater 水表手动抄表
func (h *MeterHandler) ManualReadingWater(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req ManualReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var readingTime *time.Time
	if req.ReadingTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.ReadingTime); err == nil {
			readingTime = &t
		}
	}

	// 获取操作员信息
	var operatorID *int64
	var operatorName *string
	if userID, exists := c.Get("userID"); exists {
		uid := userID.(int64)
		operatorID = &uid
		if username := middleware.GetUsername(c); username != "" {
			operatorName = &username
		}
	}

	input := &service.ManualReadingInput{
		MeterID:      id,
		MeterType:    models.MeterTypeWater,
		ReadingValue: req.ReadingValue,
		ReadingTime:  readingTime,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Remark:       req.Remark,
	}

	if err := h.meterService.ManualReading(input); err != nil {
		if err == service.ErrMeterNotFound {
			response.NotFound(c, "水表不存在")
			return
		}
		if err == service.ErrReadingTooLow {
			response.Error(c, 400, "读数不能小于当前读数")
			return
		}
		response.ServerError(c, "抄表失败")
		return
	}

	response.SuccessWithMessage(c, "抄表成功", nil)
}

// GetElectricReadings 获取电表读数历史
func (h *MeterHandler) GetElectricReadings(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, pageSize := utils.GetPage(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	readings, total, err := h.meterService.GetReadings(id, models.MeterTypeElectric, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, readings, total, page, pageSize)
}

// GetWaterReadings 获取水表读数历史
func (h *MeterHandler) GetWaterReadings(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, pageSize := utils.GetPage(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	readings, total, err := h.meterService.GetReadings(id, models.MeterTypeWater, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, readings, total, page, pageSize)
}

// ListAllReadings 获取所有手工抄表记录列表
func (h *MeterHandler) ListAllReadings(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	var meterType *int16
	if mt := c.Query("meter_type"); mt != "" {
		v := utils.ParseInt16Query(c, "meter_type")
		if v != nil {
			meterType = v
		}
	}

	meterNo := c.Query("meter_no")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	readings, total, err := h.meterService.ListManualReadings(meterType, meterNo, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	// 只返回正常状态的记录
	var filtered []service.ManualReadingInfo
	for _, r := range readings {
		if r.Status == 1 {
			filtered = append(filtered, r)
		}
	}

	response.Page(c, filtered, total, page, pageSize)
}
