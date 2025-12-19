package bookservice

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/syralon/coconut/example/ent/book"
	"github.com/syralon/coconut/example/ent/predicate"
	"github.com/syralon/coconut/example/internal/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
	"github.com/syralon/coconut/toolkit/xslices"
)

type listController struct {
	*Dependency
}

func (c *listController) List(ctx context.Context, request *example.ListBookRequest) (*example.ListBookResponse, error) {
	query := c.client.Query().Where(
		field.Selectors[predicate.Book](
			request.GetId().Selector(book.FieldID),
			request.GetTitle().Selector(book.FieldTitle),
		)...,
	)

	if paginator := request.GetPaginator(); paginator != nil {
		switch page := paginator.GetPaginator().(type) {
		case *field.Paginator_Classical:
			total, err := query.Count(ctx)
			if err != nil {
				return nil, err
			}
			for _, order := range request.GetOrders() {
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
		return nil, err
	}
	return &example.ListBookResponse{
		Data:      xslices.Trans(data, message.BookToProto),
		Paginator: request.GetPaginator(),
	}, nil
}
