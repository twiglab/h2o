package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

/*
id , uuidv7
code 充值的代码（device_code + 时间到秒）
deive_code
device_type
pos_code , 商铺号
project， 1006
充值归属（充值人）

top （限额，这里指度数，冲到多少度， = 当前表显+充电度数）
当前度数（充电时候的表显）
充电度数（冲了多少度， 例如冲50度电）

充值金额（ 反算值 电费度数差值 / （单价+ 服务费））
电费单价
用电服务费单价
充值时间
status
is_del
*/

type ValueCharge struct {
	ent.Schema
}

func (ValueCharge) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().DefaultFunc(id).SchemaType(char(36)),

		field.String("code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("充值编号"),
		// field.String("no").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("充值单号"),

		field.String("device_code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备号"),
		field.String("device_type").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("设备类型"),

		field.String("pos_code").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("位置编号"),
		field.String("project").Immutable().NotEmpty().SchemaType(varchar(64)).Comment("项目编号"),

		field.String("owner").Immutable().Optional().SchemaType(varchar(64)).Comment("归属"),

		field.Int64("top").Immutable().Default(0).Comment("限额"),     // top （限额，这里指度数，冲到多少度， = 当前表显+充电度数）
		field.Int64("stock").Immutable().Default(0).Comment("当前数值"), // 当前度数（充电时候的表显）
		field.Int64("incr").Immutable().Default(0).Comment("充值数量"),  // 充电度数（冲了多少度， 例如冲50度电）

		field.Int64("amount").Immutable().Default(0).Comment("充值金额"),   // 充值金额（ 反算值 电费度数差值 / （单价+ 服务费））
		field.Int64("unit_price").Immutable().Default(0).Comment("单价"), // 单价

		field.Time("charge_time").Immutable().Comment("充值时间"),
		field.String("charge_ts").Immutable().NotEmpty().SchemaType(varchar(36)).Comment("充值时间字符串"),

		field.Int("alarm1").Default(0).Comment("报警1"),
		field.Int("alarm2").Default(0).Comment("报警2"),
		field.Int("alarm3").Default(0).Comment("报警3"),

		field.Int("status").Default(0).Comment("状态"), // 备用

		field.String("memo").Optional().SchemaType(varchar(128)).Comment("备注"),
		field.Int("is_del").Default(0).Comment("软删除"),
	}
}

func (ValueCharge) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (ValueCharge) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("device_code"),
		index.Fields("device_type"),

		index.Fields("code").Unique(),
		index.Fields("charge_time"),
		index.Fields("charge_ts"),

		// index.Fields("pos_code"),

		index.Fields("is_del"),
	}
}

func (ValueCharge) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "value_charge"},
	}
}
