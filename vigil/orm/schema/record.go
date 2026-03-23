package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

func cdrid() string {
	u, _ := uuid.NewV7()
	return u.String()
}

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

type Record struct {
	ent.Schema
}

func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().DefaultFunc(cdrid).SchemaType(char(36)),

		field.String("device_code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
		field.String("device_type").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),

		field.String("data_code").Immutable().Unique().NotEmpty().SchemaType(varchar(64)).Comment("当前datacode"),
		field.Int64("data_value").Immutable().Default(0).Comment("当前读数"),
		field.Time("data_time").Immutable().Comment("当前时间"),

		field.String("pos_code").Immutable().SchemaType(varchar(64)).Comment("位置编号"),
		field.String("project").Immutable().SchemaType(varchar(64)).Comment("项目编号"),
	}
}

func (Record) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Record) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("device_code"),
		index.Fields("device_type"),
		index.Fields("data_code").Unique(),
		index.Fields("data_time"),
		index.Fields("pos_code"),
		index.Fields("project"),
	}
}

func (Record) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_nh_record"},
	}
}
