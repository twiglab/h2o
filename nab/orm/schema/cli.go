package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Cli struct {
	ent.Schema
}

func (Cli) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),

		field.String("code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
		field.String("typ").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),

		field.String("url").Immutable().Optional().SchemaType(varchar(64)).Comment("设备序列号"),

		field.Uint("speed").Immutable().Default(0).Comment("当前表显"),
		field.Uint("data_bits").Immutable().Default(0).Comment("当前表显"),
		field.Uint("parity").Immutable().Default(0).Comment("当前表显"),
		field.Uint("stop_bits").Immutable().Default(0).Comment("当前表显"),

		field.Int64("timeout").Immutable().Default(0).Comment("当前表显"),

		field.Uint("endian").Immutable().Default(0).Comment("当前表显"),
		field.Uint("word_order").Immutable().Default(0).Comment("当前表显"),

		field.String("memo").Immutable().SchemaType(varchar(64)).Comment("设备序列号"),
	}
}

func (Cli) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (Cli) Indexes() []ent.Index {
	return []ent.Index{}
}

func (Cli) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cli"},
	}
}
