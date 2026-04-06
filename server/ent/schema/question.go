package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Question struct {
	ent.Schema
}

func (Question) Fields() []ent.Field {
	return []ent.Field{
		field.String("question_type"),
		field.JSON("translations", map[string]string{}), // Handles JSONB translations
		field.Bool("is_active").Default(true),
	}
}

func (Question) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("answers", Answer.Type),
	}
}
