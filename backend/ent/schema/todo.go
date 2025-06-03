package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema"
    "entgo.io/contrib/entproto"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
    "time"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.Text("text").NotEmpty().
		Annotations(entproto.Field(2)),
		
		field.Time("created_at").Default(time.Now).Immutable().
		Annotations(entproto.Field(3)),
		
		field.Enum("status").NamedValues(
			"InProgress", "IN_PROGRESS",
			"Completed", "COMPLETED",
		).Default("IN_PROGRESS").
		Annotations(
			entproto.Field(4),
			entproto.Enum(map[string]int32{
				"IN_PROGRESS": 0,
				"COMPLETED": 1,
			}),
		),
		
		field.Int("priority").Default(0).
		Annotations(entproto.Field(5)),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		// 子から親
		edge.To("parent", Todo.Type).
			Unique(). // 親は単一
			Annotations(entproto.Field(102)).
			From("children").
			Annotations(entproto.Field(103)),
	}
}

func (Todo) Annotations() []schema.Annotation {
	return []schema.Annotation{
        entproto.Message(),
	}
}
