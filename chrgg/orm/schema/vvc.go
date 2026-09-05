package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type VVC struct {
	ent.Schema
}

func (VVC) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().SchemaType(char(36)),

		field.String("device_code").Immutable().SchemaType(varchar(64)).Comment("设备号"),
		field.String("device_type").Immutable().SchemaType(varchar(64)).Comment("设备类型"),

		// 这里都是指计量值（电量）
		field.Int64("quota").Immutable().Comment("限量"), // 可以用到多少度
		// field.Int64("stock").Immutable().Comment("存量"), // 当前数值（存量）
		// field.Int64("incr").Immutable().Comment("增量"),  // 增加多少（数值上等于 限量 - 存量）
		// ------

		field.Int("status").Immutable().Comment("状态"),
	}
}

func (VVC) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (VVC) Indexes() []ent.Index {
	return []ent.Index{}
}

func (VVC) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "v_vc"},
	}
}
