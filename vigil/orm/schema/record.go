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

type NhRecord struct {
	ent.Schema
}

func (NhRecord) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().DefaultFunc(cdrid).SchemaType(char(36)),

		field.String("p_code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备位置业务编号"),

		field.String("device_sn").Immutable().Optional().SchemaType(varchar(64)).Comment("设备序列号"),
		field.String("device_code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
		field.String("device_type").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),
		field.String("device_name").Immutable().Optional().SchemaType(varchar(64)).Comment("设备名称"),

		field.Int64("data_value").Immutable().Default(0).Comment("当前表显"),

		field.String("data_code").Immutable().Unique().NotEmpty().SchemaType(varchar(64)).Comment("当前记录code"),
		field.Time("data_time").Immutable().Comment("采集时间"),
		field.String("data_ts").Immutable().NotEmpty().SchemaType(varchar(36)).Comment("采集时间字符串"),

		field.String("pos_code").Immutable().SchemaType(varchar(64)).Comment("位置编号"),
		field.String("project").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("项目编号"),
		field.String("owner").Immutable().Optional().SchemaType(varchar(64)).Comment("归属"),
	}
}

func (NhRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (NhRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("p_code"),
		index.Fields("device_code"),
		index.Fields("device_type"),
		index.Fields("data_code").Unique(),
		index.Fields("data_time"),
		index.Fields("data_ts"),
		index.Fields("pos_code"),
		index.Fields("project"),
		index.Fields("owner"),
	}
}

func (NhRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "nh_record"},
	}
}
