package nab

import (
	"context"

	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/nab/orm/ent/dev"
)

type Global struct {
	Project string
	BoxCode string

	Client *ent.Client

	clientMap map[string]*ModbusCli
}

func (g Global) GetClient(code string) (*ModbusCli, bool) {
	cli, ok := g.clientMap[code]
	return cli, ok
}

func (g Global) MustGetClient(code string) *ModbusCli {
	cli, ok := g.GetClient(code)
	if !ok {
		panic("no cli " + code)
	}
	return cli
}

func (g *Global) BuildClient(ctx context.Context) error {
	g.clientMap = make(map[string]*ModbusCli)

	q := g.Client.Cli.Query()

	clients, err := q.All(ctx)
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

func (g Global) MustGetDev(ctx context.Context, code string) *ent.Dev {
	dev, err := g.GetDev(ctx, code)
	if err != nil {
		panic(err)
	}
	return dev
}
func (g Global) GetDev(ctx context.Context, code string) (*ent.Dev, error) {
	q := g.Client.Dev.Query()
	q.Where(dev.CodeEQ(code))
	d, err := q.Only(ctx)
	return d, err
}
