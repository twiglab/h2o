package idb

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
	"modernc.org/sqlite/vtab"

	"github.com/jmoiron/sqlx"
)

type IDB struct {
	dbx *sqlx.DB
	db  *sql.DB

	devFile string
	cliFile string
}

func NewIDB(dev, cli string) (*IDB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	if err := vtab.RegisterModule(db, "csv", &IDBModule{}); err != nil {
		return nil, err
	}
	dbx := sqlx.NewDb(db, "sqlite3")

	return &IDB{
		dbx: dbx,
		db:  db,

		devFile: cmp.Or(dev, "dev.csv"),
		cliFile: cmp.Or(cli, "cli.csv"),
	}, nil
}

func (i *IDB) Close() error {
	return i.dbx.Close()
}

func (i *IDB) Init() (err error) {
	_, err = i.dbx.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE dev USING csv(filename=%q)`, i.devFile))
	if err != nil {
		return
	}

	_, err = i.dbx.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE cli USING csv(filename=%q)`, i.cliFile))
	if err != nil {
		return
	}

	return
}

func (i *IDB) SelectContext(ctx context.Context, dest any, q string, a ...any) error {
	return i.dbx.SelectContext(ctx, dest, q)
}

func (i *IDB) GetContext(ctx context.Context, dest any, q string, a ...any) error {
	return i.dbx.GetContext(ctx, dest, q)
}
