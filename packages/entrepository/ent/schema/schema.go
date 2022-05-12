package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimeStamp struct {
	mixin.Schema
}

func (TimeStamp) Fields() []ent.Field {
	return []ent.Field{
		// field.Time("created_at").Default(time.Now).SchemaType(map[string]string{
		// 	dialect.MySQL: "datetime DEFAULT CURRENT_TIMESTAMP",
		// }).Immutable().StructTag(`json:"omitempty"`),
		// field.Time("updated_at").Default(time.Now).SchemaType(map[string]string{
		// 	dialect.MySQL: "datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
		// }).UpdateDefault(time.Now).StructTag(`json:"omitempty"`),
		field.Time("created_at").Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).Immutable().StructTag(`json:"omitempty"`),
		field.Time("updated_at").Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).UpdateDefault(time.Now).StructTag(`json:"omitempty"`),
		field.Time("deleted_at").SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}).Optional().StructTag(`json:"omitempty"`),
	}
}

// var ondeleteRestrict = entsql.Annotation{
// 	OnDelete: entsql.Restrict,
// }
// var ondeleteNoAction = entsql.Annotation{
// 	OnDelete: entsql.NoAction,
// }
var ondeleteCascade = entsql.Annotation{
	OnDelete: entsql.Cascade,
}
var ondeleteSetNull = entsql.Annotation{
	OnDelete: entsql.SetNull,
}
