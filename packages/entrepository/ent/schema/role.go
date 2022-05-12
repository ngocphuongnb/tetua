package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().StructTag(`validate:"max=255"`),
		field.String("description").Optional().StructTag(`validate:"max=255"`),
		field.Bool("root").Optional(),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("permissions", Permission.Type).
			Annotations(ondeleteCascade).
			StorageKey(edge.Column("role_id"), edge.Symbol("permission_role")),
		edge.To("users", User.Type).
			Annotations(ondeleteSetNull),
	}
}

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
