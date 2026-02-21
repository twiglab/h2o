package chrgg

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

const (
	createSql = "create or replace table %s as %s "
)

type Ruler struct {
	UnitFeeFen int64
}

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

func (d *DuckDB) Get(ctx context.Context, code string) (data Ruler, ok bool, err error) {
	row := d.db.QueryRowContext(ctx, d.getQry, code)
	err = row.Scan(&data.UnitFeeFen)
	ok = (err == nil)
	return
}

func (d *DuckDB) Set(_ context.Context, _ string, _ Ruler) error { return nil }
