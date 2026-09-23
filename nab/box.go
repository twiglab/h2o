package nab

import (
	"context"
	"errors"
	"time"

	"encoding/json/v2"

	"github.com/twiglab/h2o/pkg/common"

	"github.com/twiglab/h2o/nab/orm"
)

type Global struct {
	Project string
	Box     string
}

type ClientMgr struct {
	clientMap map[string]*ModbusCli
	IDB       *orm.IDB
}

func NewClientMgr(db *orm.IDB) (*ClientMgr, error) {
	ctx := context.Background()
	clientMap := make(map[string]*ModbusCli)

	clients, err := db.AllCli(ctx)
	if err != nil {
		return nil, err
	}

	for _, cli := range clients {
		cli, err := NewModbusCli(cli)
		if err != nil {
			return nil, err
		}
		clientMap[cli.Code] = cli
	}

	return &ClientMgr{
		IDB:       db,
		clientMap: clientMap,
	}, nil
}

func (g ClientMgr) ClientByCode(code string) (*ModbusCli, error) {
	cli, ok := g.clientMap[code]
	if !ok {
		return nil, errors.New("not found code " + code)
	}
	return cli, nil
}

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
	Box  string
	Code string
	Op   string
}

func (m OnOffLite) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (o OnOffLite) Topic() string {
	return common.H2O + "/onoff/" + o.Box + "/" + o.Code + "/" + o.Op
}

type OptChange struct {
	Box       string
	Code      string
	Type      string
	DataCode  string
	Op        string
	OptStatus int64
	DataTime  time.Time
}

func (m OptChange) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (o OptChange) Topic() string {
	return common.H2O + "/opt/" + o.Box + "/" + o.Code + "/" + o.Op
}
