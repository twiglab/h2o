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
	SN   string
	Code string
	Name string
	Type string
	Room string
	Rate sql.Null[float64]
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
		data.Project = h.Project
		data.Code = code
		h.Logger.WarnContext(ctx, "get", slog.String("code", code), slog.Any("data", data), slog.Any("error", err))
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
		SN:      d.SN,
		Project: h.Project,
		Code:    d.Code,
		Name:    d.Name,
		PosCode: d.Room,
		Factor:  float64int(d.Rate),
		PCode:   pcode(d.Code, h.Project),
	}, nil
}

func float64int(f sql.Null[float64]) int {
	if f.Valid {
		return int(f.V)
	}
	return 0
}

func pcode(c, p string) string {
	return c + "@" + p
}
