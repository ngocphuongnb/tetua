package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug"),
		field.String("description").Optional().StructTag(`validate:"max=255"`),
		field.Text("content").StructTag(`validate:"required"`),
		field.Text("content_html"),
		field.Int64("view_count").Default(0),
		field.Int64("comment_count").Default(0),
		field.Int64("rating_count").Optional().Default(0),
		field.Int64("rating_total").Optional().Default(0),
		field.Bool("draft").Optional().Default(false),
		field.Bool("approved").Optional().Default(false),
		field.Int("featured_image_id").Optional(),
		field.Int("user_id").Optional(),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("posts").Field("user_id").Unique(),
		edge.From("topics", Topic.Type).Ref("posts"),
		edge.From("featured_image", File.Type).Ref("posts").Field("featured_image_id").Unique(),
		edge.To("comments", Comment.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("post_id"), edge.Symbol("comment_post")),
	}
}

func (Post) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").StorageKey("name_idx"),
		index.Fields("view_count").StorageKey("view_count_idx"),
	}
}

func (Post) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
