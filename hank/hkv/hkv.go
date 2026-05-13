package hkv

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/twiglab/h2o/hank"

	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	Code string
	Name string
	Type string
	Room string
	Rate sql.Null[int]
}

type HankDB struct {
	Project string

	DB *sqlx.DB
}

func (h *HankDB) Get(ctx context.Context, code string) (data hank.MetaData, ok bool, err error) {
	data, err = h.GetOne(ctx, code)

	if err == nil {
		ok = true
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		ok = false
		return
	}
	return
}

func (h *HankDB) Set(_ context.Context, _ string, _ hank.MetaData) (err error) { return }

func (h *HankDB) GetOne(ctx context.Context, code string) (hank.MetaData, error) {
	var d Data
	if err := h.DB.GetContext(ctx, &d, get_sql, code); err != nil {
		return hank.MetaData{}, err
	}

	return hank.MetaData{
		Project: h.Project,
		Code:    d.Code,
		Name:    d.Name,
		PosCode: d.Room,
		Factor:  d.Rate.V,
	}, nil
}

const (
	get_sql = `
		select
			device_code as code,
			device_name as name,
			device_type as type,
			room_logic_code as room,
			inductance as rate
		from
			v_energy_device
		where
			device_code = ?
	`
)
