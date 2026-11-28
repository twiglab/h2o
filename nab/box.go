package nab

import (
	"context"
	"errors"

	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
)

type Box struct {
	Project string
	Code    string

	DB *orm.IDB

	cliMap map[string]PortOperator
	devMap map[string]*ent.Dev
}

func (b *Box) loadCli(ctx context.Context) error {
	clis, err := b.DB.AllCli(ctx)
	if err != nil {
		return err
	}

	b.cliMap = make(map[string]PortOperator)
	for _, cli := range clis {
		mc := NewModbusCli(cli)
		b.cliMap[cli.Code] = mc
	}
	return nil
}

func (b *Box) loadDev(ctx context.Context) error {
	devs, err := b.DB.AllDev(ctx)
	if err != nil {
		return err
	}

	b.devMap = make(map[string]*ent.Dev)
	for _, dev := range devs {
		b.devMap[dev.Code] = dev
	}
	return nil
}

func (b *Box) Load() error {
	clear(b.cliMap)
	clear(b.devMap)

	ctx := context.Background()

	if err := b.loadCli(ctx); err != nil {
		return err
	}

	if err := b.loadDev(ctx); err != nil {
		return err
	}

	return nil
}

func (b *Box) Close() error {
	for _, cli := range b.cliMap {
		_ = cli.Close()
	}
	return b.DB.Close()
}

func (b *Box) GetCli(code string) (PortOperator, error) {
	if isNullPort(code) {
		return nullOpt, nil
	}

	cli, ok := b.cliMap[code]
	if !ok {
		return nil, errors.New("not found code " + code)
	}
	return cli, nil
}

func (b *Box) GetDev(code string) (*ent.Dev, error) {
	cli, ok := b.devMap[code]
	if !ok {
		return nil, errors.New("not found code " + code)
	}
	return cli, nil
}

func (b *Box) GetDevCliEqu[T any](devCode string) (d *ent.Dev, p PortOperator, c T, err error) {
	if d, err = b.GetDev(devCode); err != nil {
		return
	}
	if p, err = b.GetCli(d.Cli); err != nil {
		return
	}
	c, err = equlib.From[T](d.Clazz)

	return
}

func (b *Box) GroupByCli() map[string][]*ent.Dev {
	group := make(map[string][]*ent.Dev)
	for _, dev := range b.devMap {
		group[dev.Cli] = append(group[dev.Cli], dev)
	}
	return group
}
