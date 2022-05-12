package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique(),
		field.String("display_name").Optional(),
		field.String("url").Optional(),
		field.String("provider").Optional(),
		field.String("provider_id").Optional(),
		field.String("provider_username").Optional(),
		field.String("provider_avatar").Optional(),
		field.String("email").Optional(),
		field.String("password").Optional(),
		field.Text("bio").Optional(),
		field.Text("bio_html").Optional(),
		field.Bool("active").Default(true),
		field.Int("avatar_image_id").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("user_id"), edge.Symbol("post_user")),
		edge.To("files", File.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("user_id"), edge.Symbol("file_user")),
		edge.To("comments", Comment.Type).
			Annotations(ondeleteSetNull).
			StorageKey(edge.Column("user_id"), edge.Symbol("comment_user")),
		edge.From("roles", Role.Type).Ref("users"),
		edge.From("avatar_image", File.Type).Ref("user_avatars").Field("avatar_image_id").Unique(),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider", "provider_id").Unique().StorageKey("provider_provider_id_unique"),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
