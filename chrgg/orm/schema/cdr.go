package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type CDR struct {
	ent.Schema
}

func (CDR) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty(),

		field.String("device_code").Immutable().NotEmpty().Comment("设备号"),
		field.String("device_type").Immutable().NotEmpty().Comment("设备类型"),

		field.Int64("last_data_value").Immutable().Default(0).Comment("上次次读数"),
		field.Int64("data_value").Immutable().Default(0).Comment("当前读数"),

		field.String("last_data_code").Immutable().Unique().NotEmpty().Comment("上次datacode"),
		field.String("data_code").Immutable().Unique().NotEmpty().Comment("当前datacode"),

		field.Time("last_data_time").Immutable().Comment("上次时间"),
		field.Time("data_time").Immutable().Comment("当前时间"),

		field.String("ploy_id").Immutable().Comment("计费方案ID"),
		field.String("rule_id").Immutable().Comment("计费规则ID"),

		field.Int64("value").Default(0).Immutable().Comment("计量数值"),
		field.Int64("unit_fee").Default(0).Immutable().Comment("计费单价"),
		field.Int64("fee").Default(0).Immutable().Comment("当次费用(fen)"),

		field.String("pos_code").Immutable().NotEmpty().Comment("位置编号"),
		field.String("project").Immutable().NotEmpty().Comment("项目编号"),

		field.String("memo").Optional().Comment("备注"),
	}
}
func (CDR) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (CDR) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("device_code"),
		index.Fields("device_type"),
		index.Fields("data_code").Unique(),
		index.Fields("data_time"),
	}
}

func (CDR) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_cdr"},
	}
}
