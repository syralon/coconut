package data

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/ent/bookshelf"
	"github.com/syralon/coconut/example/ent/predicate"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
)

type BookShelfRepository struct {
	*ent.BookShelfClient
}

func NewBookShelfRepository(client *ent.Client) *BookShelfRepository {
	return &BookShelfRepository{BookShelfClient: client.BookShelf}
}

func (rep *BookShelfRepository) Create(ctx context.Context, create *example.BookShelfCreate) (*ent.BookShelf, error) {
	op := rep.BookShelfClient.Create().
		SetName(create.GetName())
	if len(create.GetBookIds()) > 0 {
		op.AddRelBookIDs(create.GetBookIds()...)
	}
	return op.Save(ctx)
}

func (rep *BookShelfRepository) Update(ctx context.Context, id int64, update *example.BookShelfUpdate) error {
	op := rep.UpdateOneID(id)
	field.Call(op.Mutation().SetName, update.Name)
	op.Mutation().AddRelBookIDs(update.BookIds...)
	_, err := op.Save(ctx)
	return err
}

func (rep *BookShelfRepository) List(ctx context.Context, options *example.BookShelfOptions, paginator *field.Paginator) ([]*ent.BookShelf, *field.Paginator, error) {
	query := rep.Query().Where(
		field.Selectors[predicate.BookShelf](
			options.GetId().Selector(bookshelf.FieldID),
			options.GetName().Selector(bookshelf.FieldName),
		)...,
	)
	if options.GetWithBooks() {
		query = query.WithRelBooks()
	}

	if paginator != nil {
		switch page := paginator.GetPaginator().(type) {
		case *field.Paginator_Classical:
			total, err := query.Count(ctx)
			if err != nil {
				return nil, paginator, err
			}
			page.Classical.Total = int64(total)
			query = query.Order(page.Classical.OrderSelector()).
				Offset(int(page.Classical.GetLimit() * (page.Classical.GetPage() - 1))).
				Limit(int(page.Classical.GetLimit()))
		case *field.Paginator_Infinite:
			query = query.Order(bookshelf.ByID(sql.OrderDesc())).
				Limit(int(page.Infinite.GetLimit()))
			if sequence := page.Infinite.GetSequence(); sequence > 0 {
				query = query.Where(bookshelf.IDLT(page.Infinite.GetSequence()))
			}
		}
	}
	data, err := query.All(ctx)
	if err != nil {
		return nil, paginator, err
	}
	return data, paginator, nil
}

func (rep *BookShelfRepository) Raw() *ent.BookShelfClient {
	return rep.BookShelfClient
}
