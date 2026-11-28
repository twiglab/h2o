package nab

import (
	"encoding/json/v2"
	"errors"
	"strings"

	"github.com/twiglab/h2o/proto"
)

// isNullPort 判断是否为空端口
// 空端口的作用为设备不需要client 即可正常工作
// 目前为空端口为 null，nil，zero， 空端口不区分大小写
func isNullPort(code string) bool {
	switch strings.ToLower(code) { // 保留端口不区分大小写
	case "null", "nil", "zero":
		return true
	}
	return false
}

func SubscriptTopic(box string) string {
	return proto.MakeTopic(proto.ONOFF, box, "+", "+")
}

func parseOnoffTopicPart(t string) (box string, code string, op string, err error) {
	// h2o/onoff/836/836_11/ON
	//	1   2     3    4    5
	var parts []string
	if parts, err = proto.TopicPart(t); err != nil {
		return
	}

	if len(parts) != 5 {
		err = errors.New("topic error")
		return
	}

	if parts[1] != proto.ONOFF {
		err = errors.New("not h2o onoff")
		return
	}

	box, code, op = parts[2], parts[3], parts[4]
	if box == "" || code == "" {
		err = errors.New("box or code is nil")
	}

	return
}

type OnoffLite struct {
	Box  string `json:"box"`
	Code string `json:"code"`
	Op   string `json:"op"`
}

func (o OnoffLite) Topic() string {
	return proto.MakeTopic(proto.ONOFF, o.Box, o.Code, o.Op)
}

func (d OnoffLite) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
