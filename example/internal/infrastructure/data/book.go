package data

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/ent/book"
	"github.com/syralon/coconut/example/ent/predicate"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
)

var BookOrderBy = map[example.BookOrder]string{
	example.BookOrder_BookOrderByID: book.FieldID,
}

type BookRepository struct {
	*ent.BookClient
}

func NewBookRepository(client *ent.Client) *BookRepository {
	return &BookRepository{BookClient: client.Book}
}

func (rep *BookRepository) Create(ctx context.Context, create *example.BookCreate) (*ent.Book, error) {
	return rep.BookClient.Create().SetTitle(create.GetTitle()).SetAbstract(create.GetAbstract()).Save(ctx)
}

func (rep *BookRepository) Update(ctx context.Context, id int64, update *example.BookUpdate) error {
	u := rep.UpdateOneID(id)
	field.Call(u.Mutation().SetTitle, update.Title)
	field.Call(u.Mutation().SetAbstract, update.Abstract)
	_, err := u.Save(ctx)
	return err
}

func (rep *BookRepository) List(ctx context.Context, options *example.BookOptions, paginator *field.Paginator) ([]*ent.Book, *field.Paginator, error) {
	query := rep.Query().Where(
		field.Selectors[predicate.Book](
			options.GetId().Selector(book.FieldID),
			options.GetTitle().Selector(book.FieldTitle),
		)...,
	)

	if paginator != nil {
		switch page := paginator.GetPaginator().(type) {
		case *field.Paginator_Classical:
			total, err := query.Count(ctx)
			if err != nil {
				return nil, paginator, err
			}
			for _, order := range options.GetOrders() {
				if order == nil {
					continue
				}
				var opts []sql.OrderTermOption
				if order.GetDesc() {
					opts = append(opts, sql.OrderDesc())
				}
				query = query.Order(sql.OrderByField(BookOrderBy[order.GetBy()], opts...).ToFunc())
			}
			page.Classical.Total = int64(total)
			query = query.Order(page.Classical.OrderSelector()).
				Offset(int(page.Classical.GetLimit() * (page.Classical.GetPage() - 1))).
				Limit(int(page.Classical.GetLimit()))
		case *field.Paginator_Infinite:
			query = query.Order(book.ByID()).
				Limit(int(page.Infinite.GetLimit()))
			if sequence := page.Infinite.GetSequence(); sequence > 0 {
				query = query.Where(book.IDLT(page.Infinite.GetSequence()))
			}
		}
	}
	data, err := query.All(ctx)
	if err != nil {
		return nil, paginator, err
	}
	return data, paginator, nil
}

func (rep *BookRepository) Raw() *ent.BookClient {
	return rep.BookClient
}
