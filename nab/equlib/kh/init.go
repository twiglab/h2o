package kh

import (
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/proto"
)

func init() {
	equlib.Register(KH_DTZY2851_CD7_1, DTZY{})
	equlib.Register(KH_LXS_40_MI, LXSNoValve{})
}

// 宽厚厂商设备关于开合状态定义
// 0x01 为拉闸
// 0x02 为合闸
// 其他为未知状态
func optStatus(opt uint16) int64 {
	switch opt {
	case 0x01:
		return proto.OPT_STATUS_OFF
	case 0x02:
		return proto.OPT_STATUS_ON
	default:
		return proto.OPT_STATUS_UNKNOW
	}
}
