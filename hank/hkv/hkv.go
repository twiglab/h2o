package hkv

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/twiglab/h2o/hank"

	_ "github.com/go-sql-driver/mysql"
)

const Key = "hkv"

type Data struct {
	Code string
	Name string
	Type string
	Room string
	Rate sql.Null[int]
}

type HankDBConf struct {
	Project string
	DBName  string
	DSN     string
	SQLGet  string
	Logger  *slog.Logger
}

type HankDB struct {
	Project string
	DB      *sqlx.DB
	Logger  *slog.Logger
	SQLGet  string
}

func NewHankDB(conf HankDBConf) (*HankDB, error) {
	db, err := sqlx.Connect(conf.DBName, conf.DSN)
	if err != nil {
		return nil, err
	}

	return &HankDB{
		Project: conf.Project,
		DB:      db,
		Logger:  conf.Logger,
		SQLGet:  conf.SQLGet,
	}, nil
}

func (h *HankDB) Get(ctx context.Context, code string) (data hank.MetaData, ok bool, err error) {
	if data, err = h.GetOne(ctx, code); err != nil {
		h.Logger.WarnContext(ctx, "get", slog.String("code", code), slog.Any("error", err))
		return
	}
	ok = true
	return
}

func (h *HankDB) Set(_ context.Context, _ string, _ hank.MetaData) (err error) { return }

func (h *HankDB) GetOne(ctx context.Context, code string) (hank.MetaData, error) {
	var d Data
	if err := h.DB.GetContext(ctx, &d, h.SQLGet, code); err != nil {
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
