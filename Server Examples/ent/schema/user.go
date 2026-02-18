package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("google_id").Unique(),
		field.String("email").Unique(),
		field.JSON("tags", []string{}).Default([]string{}),
		field.Time("created_at").Default(time.Now),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("led_clubs", Club.Type).Ref("leaders"),
		edge.To("answers", Answer.Type),
	}
}
