package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Topic holds the schema definition for the Topic entity.
type Topic struct {
	ent.Schema
}

// Fields of the Topic.
func (Topic) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").StructTag(`validate:"max=255"`).Unique(),
		field.String("slug").StructTag(`validate:"max=255"`).Unique(),
		field.String("description").Optional().StructTag(`validate:"max=255"`),
		field.Text("content").StructTag(`validate:"required"`),
		field.Text("content_html").StructTag(`validate:"required"`),
		field.Int("parent_id").Optional(),
	}
}

// Edges of the Topic.
func (Topic) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).
			Annotations(ondeleteSetNull),
		edge.To("children", Topic.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("parent_id"), edge.Symbol("topic_parent")),
		edge.From("parent", Topic.Type).
			Field("parent_id").
			Ref("children").
			Unique(),
	}
}

func (Topic) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (Topic) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
