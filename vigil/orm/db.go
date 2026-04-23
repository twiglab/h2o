package orm

//go:generate go tool ent generate ./schema --target ./ent --feature sql/execquery,sql/upsert,privacy,sql/lock

import (
	"cmp"
	"context"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/twiglab/h2o/vigil/orm/ent/runtime"

	"github.com/twiglab/h2o/vigil/orm/ent"
)

func OpenEntClient(name, dsn string, ops ...ent.Option) (*ent.Client, error) {
	if name == "pgx" {
		return pgx(dsn, ops...)
	}
	return ent.Open(name, dsn, ops...)
}

func pgx(dsn string, ops ...ent.Option) (*ent.Client, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	db := stdlib.OpenDBFromPool(pool)
	drv := entsql.OpenDB(dialect.Postgres, db)
	ops = append(ops, ent.Driver(drv))
	return ent.NewClient(ops...), nil
}

func RecordCmp(a, b *ent.Record) int {
	return cmp.Compare(b.DataTs, a.DataTs)
}
