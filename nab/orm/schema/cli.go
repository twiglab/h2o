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

		field.String("url").Immutable().Optional().SchemaType(varchar(64)).Comment("端口URL"),

		field.Uint("speed").Immutable().Default(0).Comment("速率BPS"),
		field.Uint("data_bits").Immutable().Default(0).Comment("数据位"),
		field.Uint("parity").Immutable().Default(0).Comment("校验"),
		field.Uint("stop_bits").Immutable().Default(0).Comment("停止位"),

		field.Int64("timeout").Immutable().Default(0).Comment("超时"),

		field.Uint("endian").Immutable().Default(0).Comment("大小端"),
		field.Uint("word_order").Immutable().Default(0).Comment("字节顺序"),

		field.String("memo").Immutable().SchemaType(varchar(64)).Comment("备注"),
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
