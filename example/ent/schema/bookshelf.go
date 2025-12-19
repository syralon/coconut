package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/syralon/coconut/example/internal/global/snowflake"
)

// BookShelf holds the schema definition for the BookShelf entity.
type BookShelf struct {
	ent.Schema
}

// Fields of the BookShelf.
func (BookShelf) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(snowflake.Next),
		field.String("name"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the BookShelf.
func (BookShelf) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rel_books", Book.Type).Ref("rel_shelves"),
	}
}
