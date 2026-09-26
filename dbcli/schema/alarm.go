package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Alarm struct {
	ent.Schema
}

func (Alarm) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().DefaultFunc(id).SchemaType(char(36)),
		field.String("top_id").Unique().NotEmpty().SchemaType(char(36)).Comment("限额id"),

		field.String("top_code").Unique().NotEmpty().SchemaType(varchar(64)).Comment("限额code"),

		/*
			field.String("device_code").NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
			field.String("device_type").NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),
			field.String("device_sn").Optional().SchemaType(varchar(64)).Comment("设备序列号"),

			field.Int64("top").Immutable().Default(0).Comment("限额"),
			field.Int64("stock").Immutable().Default(0).Comment("当前数值"),
			field.Int64("incr").Immutable().Default(0).Comment("充值数量"),

			field.String("project").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("项目编号"),
			field.String("pos_code").Optional().SchemaType(varchar(64)).Comment("位置编号"),
		*/

		field.Int("alarm_1").Default(0).Comment("报警"),
		field.Time("alarm_time_1").Optional().Nillable().Comment("报警时间"),
		field.Int64("alarm_stock_1").Default(0).Comment("报警时刻表用量"),

		field.Int("alarm_2").Default(0).Comment("报警"),
		field.Time("alarm_time_2").Optional().Nillable().Comment("报警时间"),
		field.Int64("alarm_stock_2").Default(0).Comment("报警时刻表用量"),

		field.Int("alarm_3").Default(0).Comment("报警"),
		field.Time("alarm_time_3").Optional().Nillable().Comment("报警时间"),
		field.Int64("alarm_stock_3").Default(0).Comment("报警时刻表用量"),
	}
}

func (Alarm) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (Alarm) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("top_id").Unique(),
		index.Fields("top_code").Unique(),
	}
}

func (Alarm) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "mdna_alarm"},
	}
}
