package bookshelfservice

import (
	"context"

	"github.com/syralon/coconut/example/ent/bookshelf"
	"github.com/syralon/coconut/example/ent/predicate"
	"github.com/syralon/coconut/example/internal/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
	"github.com/syralon/coconut/toolkit/xslices"
)

type listController struct {
	*Dependency
}

func (c *listController) List(ctx context.Context, request *example.ListBookShelfRequest) (*example.ListBookShelfResponse, error) {
	query := c.client.Query().Where(
		field.Selectors[predicate.BookShelf](
			request.GetId().Selector(bookshelf.FieldID),
			request.GetName().Selector(bookshelf.FieldName),
		)...,
	)
	if request.GetWithBooks() {
		query = query.WithRelBooks()
	}

	if paginator := request.GetPaginator(); paginator != nil {
		switch page := paginator.GetPaginator().(type) {
		case *field.Paginator_Classical:
			total, err := query.Count(ctx)
			if err != nil {
				return nil, err
			}
			page.Classical.Total = int64(total)
			query = query.Order(page.Classical.OrderSelector()).
				Offset(int(page.Classical.GetLimit() * (page.Classical.GetPage() - 1))).
				Limit(int(page.Classical.GetLimit()))
		case *field.Paginator_Infinite:
			query = query.Order(bookshelf.ByID()).
				Limit(int(page.Infinite.GetLimit()))
			if sequence := page.Infinite.GetSequence(); sequence > 0 {
				query = query.Where(bookshelf.IDLT(page.Infinite.GetSequence()))
			}
		}
	}
	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return &example.ListBookShelfResponse{
		Data:      xslices.Trans(data, message.BookShelfToProto),
		Paginator: request.GetPaginator(),
	}, nil
}
