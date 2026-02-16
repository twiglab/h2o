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
		field.String("id"),

		field.String("device_code").Comment("设备号"),
		field.String("device_type").Comment("设备类型"),

		field.Int64("last_data_value").Comment("最后一次读数"),
		field.Int64("data_value").Comment("当前读数"),

		field.String("last_data_code").Comment("用户ID"),
		field.String("data_code").Comment("用户ID"),

		field.Time("last_data_time").Comment("上一次时间"),
		field.Time("data_time").Comment("当前时间"),

		field.String("ploy_id").Comment("计费方案ID"),
		field.String("rule_id").Comment("计费规则ID"),

		field.Int64("value").Comment("计量数值"),
		field.Int64("unit_fee").Comment("计费单价"),
		field.Int64("fee").Comment("当次费用"),

		field.String("pos_code").Comment("位置编号"),
		field.String("project").Comment("项目编号"),

		field.Time("time").Comment("处理时间"),

		field.String("remark").Comment("备注"),
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
		index.Fields("data_code"),
		index.Fields("data_time"),
	}
}

func (CDR) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_cdr"},
	}
}
