package nab

import (
	"context"

	"github.com/twiglab/h2o/nab/orm"
)

type Global struct {
	Project string
	BoxCode string

	IDB *orm.IDB

	clientMap map[string]*ModbusCli
}

func (g Global) ClientByCode(code string) *ModbusCli {
	cli := g.clientMap[code]
	return cli
}

func (g *Global) InitClients(ctx context.Context) error {
	g.clientMap = make(map[string]*ModbusCli)

	clients := g.IDB.MustAllCli(ctx)
	for _, cr := range clients {
		cli, err := NewModbusCli(cr)
		if err != nil {
			return err
		}
		g.clientMap[cli.Code] = cli
	}

	return nil
}
