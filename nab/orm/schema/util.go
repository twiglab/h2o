package schema

import (
	"fmt"

	"entgo.io/ent/dialect"
)

func char(size int) map[string]string {
	return map[string]string{
		dialect.MySQL:    fmt.Sprintf("char(%d)", size),
		dialect.SQLite:   fmt.Sprintf("char(%d)", size),
		dialect.Postgres: fmt.Sprintf("char(%d)", size),
	}
}
func varchar(size int) map[string]string {
	return map[string]string{
		dialect.MySQL:    fmt.Sprintf("varchar(%d)", size),
		dialect.SQLite:   fmt.Sprintf("varchar(%d)", size),
		dialect.Postgres: fmt.Sprintf("varchar(%d)", size),
	}
}
