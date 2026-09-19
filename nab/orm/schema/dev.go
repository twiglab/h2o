package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Dev struct {
	ent.Schema
}

func (Dev) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),

		field.String("code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
		field.String("typ").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),
		field.String("sn").Immutable().Optional().SchemaType(varchar(64)).Comment("设备序列号"),

		field.String("clazz").Immutable().Unique().NotEmpty().SchemaType(varchar(64)).Comment("当前记录code"),

		field.String("cli").Immutable().NotEmpty().SchemaType(varchar(36)).Comment("采集时间字符串"),

		field.Uint8("unit_id").Immutable().Default(0).Comment("当前表显"),
		field.Uint("endian").Immutable().Default(0).Comment("当前表显"),
		field.Uint("word_order").Immutable().Default(0).Comment("当前表显"),
	}
}

func (Dev) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (Dev) Indexes() []ent.Index {
	return []ent.Index{}
}

func (Dev) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "dev"},
	}
}
