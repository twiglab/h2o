package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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

		field.Time("last_data_time").Comment("签到时间"),
		field.Time("data_time").Comment("签到时间"),

		field.String("ploy_id").Comment("计费方案ID"),
		field.String("rule_id").Comment("计费规则ID"),

		field.Int64("value").Comment("计费数值"),
		field.Int64("unit_fee").Comment("计费单价"),
		field.Int64("fee").Comment("当次费用"),

		field.String("remark").Default("").Comment("备注"),
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
