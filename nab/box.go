package nab

import (
	"context"

	"github.com/twiglab/h2o/nab/orm"
)

type Global struct {
	Project string
	Box     string

	IDB *orm.IDB

	clientMap map[string]*ModbusCli
}

func (g Global) ClientByCode(code string) (cli *ModbusCli, ok bool) {
	cli, ok = g.clientMap[code]
	return
}

func (g *Global) InitClients(ctx context.Context) error {
	g.clientMap = make(map[string]*ModbusCli)

	clients, err := g.IDB.AllCli(ctx)
	if err != nil {
		return err
	}

	for _, cr := range clients {
		cli, err := NewModbusCli(cr)
		if err != nil {
			return err
		}
		g.clientMap[cli.Code] = cli
	}

	return nil
}
