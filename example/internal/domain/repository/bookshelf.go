package repository

import (
	"context"

	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
)

type BookShelfRepository interface {
	Create(ctx context.Context, create *example.BookShelfCreate) (*ent.BookShelf, error)
	Update(ctx context.Context, id int64, update *example.BookShelfUpdate) error
	List(ctx context.Context, options *example.BookShelfOptions, paginator *field.Paginator) ([]*ent.BookShelf, *field.Paginator, error)
	//Raw() *ent.BookShelfClient
}
