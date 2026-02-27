package abm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/duckdb/duckdb-go/v2"
)

const (
	duckdb = "duckdb"
)

func qry(get, lst, tbl string) (g string, l string) {
	g = fmt.Sprintf(get, tbl)
	if lst != "" {
		l = fmt.Sprintf(lst, tbl)
	}
	return
}

func nextTbl(curr string) string {
	if curr == "db_a" {
		return "db_b"
	}

	return "db_a"
}

func loadSql(createSql, lastTbl string) (tbl string, cr string) {
	tbl = nextTbl(lastTbl)
	cr = fmt.Sprintf(createSql, tbl)
	return
}

type Conf struct {
	LoadSQL string
	GetSQL  string
	ListSQL string

	Period int
}

type DuckABM[K comparable, T any] struct {
	dbx *sqlx.DB

	conf Conf

	tbl     string
	getQry  string
	listQry string
}

func NewDuckABM[K comparable, T any](conf Conf) (*DuckABM[K, T], error) {
	db, err := sqlx.Connect(duckdb, "")
	if err != nil {
		return nil, err
	}

	if conf.Period == 0 {
		conf.Period = 60
	}

	return &DuckABM[K, T]{
		dbx:  db,
		conf: conf,
	}, nil
}

func (d *DuckABM[K, T]) Load(ctx context.Context) error {
	nextTbl, cr := loadSql(d.conf.LoadSQL, d.tbl)
	slog.DebugContext(ctx, "ddbLoad",
		slog.String("tbl", d.tbl),
		slog.String("nextTbl", nextTbl),
		slog.String("create", cr),
	)
	if _, err := d.dbx.ExecContext(ctx, cr); err != nil {
		return err
	}

	d.tbl = nextTbl
	d.getQry, d.listQry = qry(d.conf.GetSQL, d.conf.ListSQL, d.tbl)
	return nil
}

func (d *DuckABM[K, T]) Loop(ctx context.Context) error {
	if err := d.Load(ctx); err != nil {
		return err
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Minute * time.Duration(d.conf.Period))
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
func (d *DuckABM[K, T]) List(ctx context.Context) (ds []T, err error) {
	if d.conf.ListSQL == "" {
		return nil, errors.New("no list sql")
	}

	err = d.dbx.SelectContext(ctx, &ds, d.listQry)
	return
}

func (d *DuckABM[K, T]) Get(ctx context.Context, code K) (data T, ok bool, err error) {
	err = d.dbx.GetContext(ctx, &data, d.getQry, code)
	ok = (err == nil)
	return
}

func (d *DuckABM[K, T]) Set(_ context.Context, _ K, _ T) (err error)             { return }
func (d *DuckABM[K, T]) Clear(_ context.Context) (err error)                     { return }
func (d *DuckABM[K, T]) Forget(_ context.Context, _ K) (v T, ok bool, err error) { return }
