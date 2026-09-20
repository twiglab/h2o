package orm

//go:generate go tool ent generate ./schema --target ./ent --feature sql/execquery,sql/upsert,privacy,sql/lock

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"strconv"

	_ "modernc.org/sqlite"
	"modernc.org/sqlite/vtab"

	"github.com/twiglab/h2o/nab/orm/ent/dev"
	_ "github.com/twiglab/h2o/nab/orm/ent/runtime"

	"github.com/twiglab/h2o/nab/orm/idb"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"github.com/twiglab/h2o/nab/orm/ent"
)

type IDB struct {
	cli *ent.Client
}

func (g IDB) GetDev(ctx context.Context, code string) (*ent.Dev, error) {
	q := g.cli.Dev.Query()
	q.Where(dev.CodeEQ(code))
	d, err := q.Only(ctx)
	return d, err
}

func (g IDB) AllDev(ctx context.Context) ([]*ent.Dev, error) {
	q := g.cli.Dev.Query()
	return q.All(ctx)
}

func (g IDB) AllCli(ctx context.Context) ([]*ent.Cli, error) {
	q := g.cli.Cli.Query()
	return q.All(ctx)
}

func (g IDB) Close() error {
	return g.cli.Close()
}

func NewIDB(dev, cli string, ops ...ent.Option) (*IDB, error) {
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

	return &IDB{
		cli: ent.NewClient(ops...),
	}, nil
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
