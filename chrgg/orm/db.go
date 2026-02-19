package orm

//go:generate go tool ent generate ./schema --target ./ent --feature sql/execquery,sql/upsert,privacy,sql/lock

import (
	"context"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/twiglab/h2o/chrgg/orm/ent/runtime"

	"github.com/twiglab/h2o/chrgg/orm/ent"
)

func OpenEntClient(name, dns string) (*ent.Client, error) {
	if name == "pgx" {
		pool, err := pgxpool.New(context.Background(), dns)
		if err != nil {
			return nil, err
		}
		db := stdlib.OpenDBFromPool(pool)
		drv := entsql.OpenDB(dialect.Postgres, db)
		return ent.NewClient(ent.Driver(drv)), nil
	}

	return ent.Open(name, dns)
}
