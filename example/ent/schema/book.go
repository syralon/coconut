package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/syralon/coconut/example/internal/global/snowflake"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(snowflake.Next),
		field.String("title"),
		field.String("abstract").Optional(),
		field.String("author").Optional(),
		field.String("cover").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rel_shelves", BookShelf.Type),
	}
}
