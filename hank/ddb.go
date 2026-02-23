package hank

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/duckdb/duckdb-go/v2"
)

const (
	createSql = "create or replace table %s as %s "
)

func qry(q, tbl string) string {
	return fmt.Sprintf(q, tbl)
}

func nextTbl(curr string) string {
	if curr == "db_a" {
		return "db_b"
	}

	return "db_a"
}

func losdSql(t, load string) (tbl string, from string) {
	tbl = nextTbl(t)
	from = fmt.Sprintf(createSql, tbl, load)
	return
}

func cutWhere(s string) string {
	b, _, ok := strings.Cut(s, "where")
	if !ok {
		panic(s)
	}
	return strings.TrimSpace(b)
}

type MetaData struct {
	SN   string `json:"sn,omitempty"`   // 仪表的序列号,仪表上有个条形码,如果没有就是空,或者自定义
	Code string `json:"code"`           // 设备code,业务全局唯一
	Name string `json:"name,omitempty"` // 设备名称,可以为空

	Project   string `json:"project"`                // 所属项目编号
	PosCode   string `json:"pos_code" db:"pos_code"` // 位置编号
	Building  string `json:"building"`               // 大楼
	FloorCode string `json:"floor_code" db:"floor_code"`
	AreaCode  string `json:"area_code" db:"area_code"`

	F1 string `json:"f1"`
	F2 string `json:"f2"`
	F3 string `json:"f3"`
	F4 string `json:"f4"`
	F5 string `json:"f5"`
}

func (m MetaData) ToStrings() []string {
	return []string{
		m.SN, m.Code, m.Name, m.Project, m.PosCode, m.Building, m.FloorCode, m.AreaCode,
		m.F1, m.F2, m.F3, m.F4, m.F5,
	}
}

type DuckDB struct {
	db *sqlx.DB

	from string
	q    string

	tbl    string
	getQry string
}

func NewDDB(from, q string) (*DuckDB, error) {
	db, err := sqlx.Connect("duckdb", "")
	if err != nil {
		return nil, err
	}

	return &DuckDB{
		db:  db,
		tbl: "x",

		from: from,
		q:    q,
	}, nil
}

func (d *DuckDB) Load(ctx context.Context) error {
	nextTbl, cr := losdSql(d.tbl, d.from)
	slog.DebugContext(ctx, "ddbLoad",
		slog.String("tbl", d.tbl),
		slog.String("nextTbl", nextTbl),
		slog.String("create", cr),
	)
	if _, err := d.db.ExecContext(ctx, cr); err != nil {
		return err
	}

	d.tbl = nextTbl
	d.getQry = qry(d.q, d.tbl)
	return nil
}

func (d *DuckDB) Loop(ctx context.Context) error {
	if err := d.Load(ctx); err != nil {
		return err
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = d.Load(ctx)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return nil
}

func (d *DuckDB) List(ctx context.Context) (mds []MetaData, err error) {
	s := cutWhere(d.getQry)
	err = d.db.SelectContext(ctx, &mds, s)
	return
}

func (d *DuckDB) Get(ctx context.Context, code string) (data MetaData, ok bool, err error) {
	err = d.db.GetContext(ctx, &data, d.getQry, code)
	ok = (err == nil)
	return
}

func (d *DuckDB) Set(_ context.Context, _ string, _ MetaData) (err error)               { return }
func (d *DuckDB) Clear(_ context.Context) (err error)                                   { return }
func (d *DuckDB) Forget(_ context.Context, _ string) (val MetaData, ok bool, err error) { return }
