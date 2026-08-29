package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ViewRecord struct {
	ent.Schema
}

func (ViewRecord) Fields() []ent.Field {
	return []ent.Field{
		// nhrecord.id
		field.String("id").Immutable().SchemaType(char(36)),

		// ----------------------------------
		// record a
		field.String("device_code").Immutable().SchemaType(varchar(64)).Comment("设备号"),
		field.String("device_type").Immutable().SchemaType(varchar(64)).Comment("设备类型"),
		field.Int64("data_value").Immutable().Comment("当前表显"),
		field.Time("data_time").Immutable().Comment("采集时间"),
		field.String("data_ts").Immutable().SchemaType(varchar(36)).Comment("采集时间字符串"),

		// ----------------------------------

		// device b

		field.String("device_sn").Immutable().SchemaType(varchar(64)).Comment("设备序列号"),
		field.String("device_name").Immutable().SchemaType(varchar(64)).Comment("设备名称"),

		field.Int("rate").Immutable().Comment("倍率"),

		field.String("project").Immutable().SchemaType(varchar(64)).Comment("项目编号"),
		field.String("pos_code").Immutable().SchemaType(varchar(64)).Comment("位置编号"),

		field.Int("status").Immutable().Comment("状态"),

		field.String("memo").Immutable().SchemaType(varchar(128)).Comment("备注"),
	}
}

func (ViewRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (ViewRecord) Indexes() []ent.Index {
	return []ent.Index{}
}

func (ViewRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "view_record"},
	}
}
