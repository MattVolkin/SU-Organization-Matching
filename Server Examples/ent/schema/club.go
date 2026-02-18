package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Club struct {
	ent.Schema
}

func (Club) Fields() []ent.Field {
	return []ent.Field{
		field.String("club_name"),
		field.Text("description").Optional(),
		field.String("image_path").Optional(),
		field.String("external_link").Optional(),
		field.Text("contact_info").Optional(),
		field.Bool("include_officer_emails").Default(false),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Club) Edges() []ent.Edge {
	return []ent.Edge{
		// This creates the many-to-many relationship (club_leaders table)
		edge.To("leaders", User.Type),
	}
}
