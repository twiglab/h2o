package nab

import (
	"context"
	"encoding/json/v2"
	"errors"
	"strings"

	"github.com/twiglab/h2o/proto"

	"github.com/twiglab/h2o/nab/orm"
)

type Global struct {
	Project string
	Box     string
}

type ClientMgr struct {
	clientMap map[string]PortOperator
	IDB       *orm.IDB
}

func NewClientMgr(db *orm.IDB) (*ClientMgr, error) {
	ctx := context.Background()
	clientMap := make(map[string]PortOperator)

	clients, err := db.AllCli(ctx)
	if err != nil {
		return nil, err
	}

	for _, cli := range clients {
		cli := NewModbusCli(cli)
		if err := cli.Open(); err != nil {
			return nil, err
		}
		clientMap[cli.Code] = cli
	}

	return &ClientMgr{
		IDB:       db,
		clientMap: clientMap,
	}, nil
}

// isNullPort 判断是否为空端口
// 空端口的作用为设备不需要client 即可正常工作
// 目前为空端口为 null，nil，zero， 空端口不区分大小写
func isNullPort(code string) bool {
	switch strings.ToLower(code) {
	case "null", "nil", "zero":
		return true
	}
	return false
}

func (g ClientMgr) ClientByCode(code string) (PortOperator, error) {
	if isNullPort(code) {
		return nullOpt, nil
	}

	cli, ok := g.clientMap[code]
	if !ok {
		return nil, errors.New("not found code " + code)
	}
	return cli, nil
}

func SubscriptTopic(box string) string {
	return proto.MakeTopic(proto.ONOFF, box, "#")
}

func onoffTopicPart(t string) (box string, code string, op string, err error) {
	// h2o/onoff/836/836_11/ON
	//  1   2     3    4    5
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
