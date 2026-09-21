package nab

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/pkg/common"
)

func SubscriptTopic(boxCode string) string {
	return common.H2O + "/onoff/" + boxCode + "/#"
}

func topicPart(t string) (string, string, string) {
	ss := common.TopicPart(t)
	_ = ss[4]

	if ss[1] != "onoff" {
		panic("not h2o onoff")
	}

	return ss[2], ss[3], ss[4]
}

type OnOffLite struct {
	BoxCode string
	Code    string
	Op      string
}

func (m OnOffLite) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (o OnOffLite) Topic() string {
	return common.H2O + "/onoff/" + o.BoxCode + "/" + o.Code + "/" + o.Op
}

type OnOffer interface {
	DeviceOn
	DeviceOff
}
