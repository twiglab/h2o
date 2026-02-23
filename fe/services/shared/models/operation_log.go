package models

import "time"

// 操作日志状态
const (
	OperationLogStatusFailed  = 0 // 失败
	OperationLogStatusSuccess = 1 // 成功
)

// 操作类型
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionQuery  = "query"
	ActionExport = "export"
	ActionImport = "import"
	ActionLogin  = "login"
	ActionLogout = "logout"
)

// OperationLog 操作日志表
type OperationLog struct {
	ID            int64     `gorm:"column:id;primaryKey" json:"id"`
	UserID        *int64    `gorm:"column:user_id" json:"user_id"`
	Username      *string   `gorm:"column:username" json:"username"`
	Module        *string   `gorm:"column:module" json:"module"`
	Action        *string   `gorm:"column:action" json:"action"`
	Description   *string   `gorm:"column:description" json:"description"`
	TargetType    *string   `gorm:"column:target_type" json:"target_type"`
	TargetID      *string   `gorm:"column:target_id" json:"target_id"`
	RequestMethod *string   `gorm:"column:request_method" json:"request_method"`
	RequestURL    *string   `gorm:"column:request_url" json:"request_url"`
	RequestParams *string   `gorm:"column:request_params" json:"request_params"`
	ResponseCode  *int      `gorm:"column:response_code" json:"response_code"`
	OldValue      *string   `gorm:"column:old_value" json:"old_value"`
	NewValue      *string   `gorm:"column:new_value" json:"new_value"`
	IPAddress     *string   `gorm:"column:ip_address" json:"ip_address"`
	UserAgent     *string   `gorm:"column:user_agent" json:"user_agent"`
	DurationMs    *int      `gorm:"column:duration_ms" json:"duration_ms"`
	Status        int8      `gorm:"column:status" json:"status"`
	ErrorMsg      *string   `gorm:"column:error_msg" json:"error_msg"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (OperationLog) TableName() string { return "sys_operation_log" }

// IsSuccess 是否成功
func (l *OperationLog) IsSuccess() bool {
	return l.Status == OperationLogStatusSuccess
}

// OperationLogCreateReq 创建操作日志请求
type OperationLogCreateReq struct {
	UserID        int64
	Username      string
	Module        string
	Action        string
	Description   string
	TargetType    string
	TargetID      string
	RequestMethod string
	RequestURL    string
	RequestParams string
	ResponseCode  int
	OldValue      string
	NewValue      string
	IPAddress     string
	UserAgent     string
	DurationMs    int
	Status        int8
	ErrorMsg      string
}

// ToModel 转换为模型
func (r *OperationLogCreateReq) ToModel() *OperationLog {
	log := &OperationLog{
		Status:    r.Status,
		CreatedAt: time.Now(),
	}

	if r.UserID > 0 {
		log.UserID = &r.UserID
	}
	if r.Username != "" {
		log.Username = &r.Username
	}
	if r.Module != "" {
		log.Module = &r.Module
	}
	if r.Action != "" {
		log.Action = &r.Action
	}
	if r.Description != "" {
		log.Description = &r.Description
	}
	if r.TargetType != "" {
		log.TargetType = &r.TargetType
	}
	if r.TargetID != "" {
		log.TargetID = &r.TargetID
	}
	if r.RequestMethod != "" {
		log.RequestMethod = &r.RequestMethod
	}
	if r.RequestURL != "" {
		log.RequestURL = &r.RequestURL
	}
	if r.RequestParams != "" {
		log.RequestParams = &r.RequestParams
	}
	if r.ResponseCode > 0 {
		log.ResponseCode = &r.ResponseCode
	}
	if r.OldValue != "" {
		log.OldValue = &r.OldValue
	}
	if r.NewValue != "" {
		log.NewValue = &r.NewValue
	}
	if r.IPAddress != "" {
		log.IPAddress = &r.IPAddress
	}
	if r.UserAgent != "" {
		log.UserAgent = &r.UserAgent
	}
	if r.DurationMs > 0 {
		log.DurationMs = &r.DurationMs
	}
	if r.ErrorMsg != "" {
		log.ErrorMsg = &r.ErrorMsg
	}

	return log
}
