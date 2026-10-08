package proto

// 点位信息
type Pos struct {
	Project string         `json:"project,omitempty"`  // 所属项目编号
	PosCode string         `json:"pos_code,omitempty"` // 位置编号
	Extra   map[string]any `json:",embed"`
}
