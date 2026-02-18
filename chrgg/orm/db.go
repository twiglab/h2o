package orm

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	_ "github.com/twiglab/h2o/chrgg/orm/ent/runtime"
)

//go:generate go tool ent generate ./schema --target ./ent

func OpenEntClient(name, dns string) (*ent.Client, error) {
	return ent.Open(name, dns)
}
