package orm

//go:generate go tool ent generate ./schema --target ./ent --feature sql/execquery,sql/upsert,privacy,sql/lock

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
	"modernc.org/sqlite/vtab"

	_ "github.com/twiglab/h2o/nab/orm/ent/runtime"

	"github.com/twiglab/h2o/nab/orm/idb"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"github.com/twiglab/h2o/nab/orm/ent"
)

func Sqlite(dev, cli string, ops ...ent.Option) (*ent.Client, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	if err := vtab.RegisterModule(db, "csv", &idb.IDBModule{}); err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE dev USING csv(filename=%q)`, dev))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE cli USING csv(filename=%q)`, cli))
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	ops = append(ops, ent.Driver(drv))
	return ent.NewClient(ops...), nil
}
