package hank

import "context"

type MetaData struct {
	SN   string `json:"sn,omitempty"`   // 仪表的序列号,仪表上有个条形码,如果没有就是空,或者自定义
	Code string `json:"code"`           // 设备code,业务全局唯一
	Name string `json:"name,omitempty"` // 设备名称,可以为空

	Project string `json:"project"`            // 所属项目编号
	PosCode string `json:"pos_code,omitempty"` // 位置编号
	Owner   string `json:"owner,omitempty"`
}

type SimpleMD struct {
	Project string
}

func (s SimpleMD) Get(ctx context.Context, code string) (MetaData, bool, error) {
	return MetaData{
		Code:    code,
		Project: s.Project,
	}, true, nil
}

func (e SimpleMD) Set(_ context.Context, _ string, _ MetaData) (err error) { return }
