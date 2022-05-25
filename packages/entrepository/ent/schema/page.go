package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Page holds the schema definition for the Page entity.
type Page struct {
	ent.Schema
}

// Fields of the Page.
func (Page) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("slug"),
		field.Text("content").StructTag(`validate:"required"`),
		field.Text("content_html"),
		field.Bool("draft").Optional().Default(false),
		field.Int("featured_image_id").Optional(),
	}
}

// Edges of the Page.
func (Page) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("featured_image", File.Type).Ref("pages").Field("featured_image_id").Unique(),
	}
}

func (Page) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").StorageKey("name_idx"),
		index.Fields("slug").StorageKey("slug_unique").Unique(),
	}
}

func (Page) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}

func (Page) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeStamp{},
	}
}
