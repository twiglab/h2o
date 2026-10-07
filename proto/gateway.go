package proto

const (
	GATEWAY_NH = "NH"
)

type Gateway struct {
	Code string `json:"code"`           // 网关code,业务全局唯一
	Type string `json:"type,omitempty"` // 网关类型
	Pos  string `json:"pos,omitempty"`  // 网关所在位置
}

type Modbus struct {
	Gateway
	UnitID uint8 `json:"unit_id"`
}
