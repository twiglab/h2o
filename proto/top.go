package proto

import (
	"time"
)

const (
	STATUS_END     = 1
	STATUS_BEGIN   = 0
	STATUS_INVALID = -1
)

type Top struct {
	Code string `json:"code"`

	Top   int64 `json:"top"`
	Stock int64 `json:"stock"`
	Incr  int64 `json:"incr"`

	ChargeTime time.Time `json:"charge_time"`

	Alarm int `json:"alarm"`

	Status int `json:"status"`
}
