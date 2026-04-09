package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Answer struct {
	ent.Schema
}

func (Answer) Fields() []ent.Field {
	return []ent.Field{
		field.Int("question_id"),
		field.Int("user_id"),
		field.Text("answer_text"),
		field.Time("submitted_at").Default(time.Now),
	}
}

func (Answer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("question", Question.Type).Ref("answers").Field("question_id").Unique().Required(),
		edge.From("user", User.Type).Ref("answers").Field("user_id").Unique().Required(),
	}
}
