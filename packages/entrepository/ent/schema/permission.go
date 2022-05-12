package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"),
		field.String("action").StructTag(`validate:"max=255"`),
		field.String("value").Comment("all | own | none"),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).Ref("permissions").Unique().Field("role_id").Required(),
	}
}

func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "action").Unique().StorageKey("role_action_unique_idx"),
	}
}

func (Permission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
