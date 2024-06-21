package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type SomeTable struct {
	ent.Schema
}

// Fields of the User.
func (SomeTable) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Positive().Annotations(),
		field.String("name").Default("").Comment("").Annotations(),
		field.Uint64("created_at").Default(0).Annotations(),
		field.Uint64("updated_at").Default(0).Annotations(),
		field.Uint64("deleted_at").Default(0).Annotations(),
	}
}
