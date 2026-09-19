package orm

//go:generate go tool ent generate ./schema --target ./ent --feature sql/execquery,sql/upsert,privacy,sql/lock

import (
	"cmp"
	"database/sql"
	"fmt"
	"strconv"

	_ "modernc.org/sqlite"
	"modernc.org/sqlite/vtab"

	_ "github.com/twiglab/h2o/nab/orm/ent/runtime"

	"github.com/twiglab/h2o/nab/orm/idb"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"github.com/twiglab/h2o/nab/orm/ent"
)

func NewIDB(dev, cli string, ops ...ent.Option) (*ent.Client, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	if err := vtab.RegisterModule(db, "csv", &idb.IDBModule{}); err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE dev USING csv(filename=%q)`, cmp.Or(dev, "dev.csv")))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE cli USING csv(filename=%q)`, cmp.Or(cli, "cli.csv")))
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	ops = append(ops, ent.Driver(drv))
	return ent.NewClient(ops...), nil
}

func DevToStrings(v *ent.Dev) []string {
	return []string{
		strconv.FormatInt(v.ID, 10),
		v.Code,
		v.Typ,
		v.Sn,
		v.Clazz,
		v.Cli,

		strconv.FormatInt(int64(v.UnitID), 10),
		strconv.FormatInt(int64(v.Addr), 10),
		strconv.FormatInt(int64(v.Endian), 10),
		strconv.FormatInt(int64(v.WordOrder), 10),

		v.Memo,
	}
}

func CliToStrings(v *ent.Cli) []string {
	return []string{
		strconv.FormatInt(v.ID, 10),
		v.Code,
		v.Typ,
		v.URL,

		strconv.FormatInt(int64(v.Speed), 10),
		strconv.FormatInt(int64(v.DataBits), 10),
		strconv.FormatInt(int64(v.Parity), 10),
		strconv.FormatInt(int64(v.StopBits), 10),
		strconv.FormatInt(int64(v.Timeout), 10),

		strconv.FormatInt(int64(v.Endian), 10),
		strconv.FormatInt(int64(v.WordOrder), 10),

		v.Memo,
	}
}
