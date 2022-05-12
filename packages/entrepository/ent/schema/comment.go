package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Comment holds the schema definition for the Comment entity.
type Comment struct {
	ent.Schema
}

// Fields of the Comment.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Text("content").StructTag(`validate:"required"`),
		field.Text("content_html").StructTag(`validate:"required"`),
		field.Int64("votes").Default(0),
		field.Int("post_id").Optional(),
		field.Int("user_id").Optional(),
		field.Int("parent_id").Optional(),
	}
}

// Edges of the Comment.
func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).Ref("comments").Field("post_id").Unique(),
		edge.From("user", User.Type).Ref("comments").Field("user_id").Unique(),
		edge.To("children", Comment.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("parent_id"), edge.Symbol("comment_parent")),
		edge.From("parent", Comment.Type).
			Field("parent_id").
			Ref("children").
			Unique(),
	}
}

func (Comment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (Comment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
