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

func id() string {
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

type Device struct {
	ent.Schema
}

func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().DefaultFunc(id).SchemaType(char(36)),

		field.String("device_code").NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
		field.String("device_type").NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),

		field.String("device_sn").Optional().SchemaType(varchar(64)).Comment("设备序列号"),
		field.String("device_name").Optional().SchemaType(varchar(64)).Comment("设备名称"),

		field.Int("rate").Default(1).Comment("当前倍率"),

		field.String("project").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("项目编号"),
		field.String("pos_code").Optional().SchemaType(varchar(64)).Comment("位置编号"),
		field.String("area_code").Optional().SchemaType(varchar(64)).Comment("区域编号"),
		field.String("pcode").Optional().SchemaType(varchar(64)).Comment("对外位置编号"),

		field.Int("status").Default(0).Comment("状态"),

		field.String("memo").Optional().SchemaType(varchar(128)).Comment("备注"),

		field.Int("is_del").Default(0).Comment("软删除"),
	}
}

func (Device) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Device) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("device_code").Unique(),
		index.Fields("device_type"),
		index.Fields("device_sn"),

		index.Fields("pos_code"),
		index.Fields("project"),
		index.Fields("pcode"),
	}
}

func (Device) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "device"},
	}
}
