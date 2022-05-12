package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("disk"),
		field.String("path").SchemaType(map[string]string{
			dialect.MySQL: "varchar(500)",
		}),
		field.String("type"),
		field.Int("size"),
		field.Int("user_id").Optional(),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("files").Field("user_id").Unique(),

		edge.To("posts", Post.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("featured_image_id"), edge.Symbol("post_featured_image")),

		edge.To("user_avatars", User.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("avatar_image_id"), edge.Symbol("user_avatar_image")),
	}
}

func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("path").StorageKey("path_idx"),
	}
}

func (File) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
