package repository

import (
	"context"

	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
)

type BookRepository interface {
	Create(ctx context.Context, create *example.BookCreate) (*ent.Book, error)
	Update(ctx context.Context, id int64, update *example.BookUpdate) error
	List(ctx context.Context, options *example.BookOptions, paginator *field.Paginator) ([]*ent.Book, *field.Paginator, error)
	//Raw() *ent.BookClient
}
