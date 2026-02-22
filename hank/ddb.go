package hank

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

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

	Project   string `json:"project"`  // 所属项目编号
	PosCode   string `json:"pos_code"` // 位置编号
	Building  string `json:"building"` // 大楼
	FloorCode string `json:"floor_code"`
	AreaCode  string `json:"area_code"`

	F1 string `json:"f1"`
	F2 string `json:"f2"`
	F3 string `json:"f3"`
	F4 string `json:"f4"`
	F5 string `json:"f5"`
}

type DuckDB struct {
	db *sql.DB

	from string
	q    string

	tbl    string
	getQry string
}

func NewDDB(from, q string) (*DuckDB, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
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

func (d *DuckDB) List(ctx context.Context) ([]MetaData, error) {
	s := cutWhere(d.getQry)

	rows, err := d.db.QueryContext(ctx, s)
	if err != nil {
		return nil, err
	}
	var mds []MetaData
	for rows.Next() {
		var data MetaData
		err = rows.Scan(
			&data.SN, &data.Code, &data.Name, &data.Project,
			&data.PosCode, &data.Building, &data.FloorCode, &data.AreaCode,
			&data.F1, &data.F2, &data.F3, &data.F4, &data.F5,
		)
		if err != nil {
			return mds, err
		}
		mds = append(mds, data)
	}

	return mds, nil
}

func (d *DuckDB) Get(ctx context.Context, code string) (data MetaData, ok bool, err error) {
	row := d.db.QueryRowContext(ctx, d.getQry, code)
	err = row.Scan(
		&data.SN, &data.Code, &data.Name, &data.Project,
		&data.PosCode, &data.Building, &data.FloorCode, &data.AreaCode,
		&data.F1, &data.F2, &data.F3, &data.F4, &data.F5,
	)
	ok = err == nil
	return
}

func (d *DuckDB) Set(_ context.Context, _ string, _ MetaData) error { return nil }
