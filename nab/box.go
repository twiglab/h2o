package nab

import (
	"context"
	"errors"

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
