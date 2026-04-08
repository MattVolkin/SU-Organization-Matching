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
		field.String("meeting_time").Optional(),
		field.String("image_path").Optional(),
		field.String("external_link").Optional(),
		field.Text("contact_info").Optional(),
		field.Bool("include_officer_emails").Default(false),
		field.JSON("personality", []string{}).Default([]string{}),
		field.JSON("activities", []string{}).Default([]string{}),
		field.JSON("genders", []string{}).Default([]string{}),
		field.JSON("ethnicities", []string{}).Default([]string{}),
		field.JSON("religions", []string{}).Default([]string{}),
		field.Bool("strict_genders").Default(false),
		field.JSON("dedicated_majors", []string{}).Default([]string{}),
		field.JSON("associated_majors", []string{}).Default([]string{}),
		field.JSON("other", []string{}).Default([]string{}),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Club) Edges() []ent.Edge {
	return []ent.Edge{
		// This creates the many-to-many relationship (club_leaders table)
		edge.To("leaders", User.Type),
	}
}
