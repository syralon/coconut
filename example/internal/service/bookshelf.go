package service

import (
	"context"

	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/ent/bookshelf"
	"github.com/syralon/coconut/example/ent/predicate"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/proto/syralon/coconut/field"
	"github.com/syralon/coconut/toolkit/xslices"
)

type BookShelfService struct {
	example.UnimplementedBookShelfServiceServer

	client *ent.BookShelfClient
}

func NewBookShelfService(client *ent.Client) *BookShelfService {
	return &BookShelfService{client: client.BookShelf}
}

func BookShelfToProto(data *ent.BookShelf) *example.BookShelf {
	shelf := &example.BookShelf{
		Name: data.Name,
	}
	if data.Edges.Books != nil {
		shelf.Books = make([]*example.Book, 0, len(data.Edges.Books))
	}
	for _, book := range data.Edges.Books {
		shelf.Books = append(shelf.Books, BookToProto(book))
	}
	return shelf
}

func (s *BookShelfService) Create(ctx context.Context, request *example.CreateBookShelfRequest) (*example.CreateBookShelfResponse, error) {
	data, err := s.client.Create().
		SetName(request.GetName()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &example.CreateBookShelfResponse{Data: BookShelfToProto(data)}, nil
}

func (s *BookShelfService) List(ctx context.Context, request *example.ListBookShelfRequest) (*example.ListBookShelfResponse, error) {
	query := s.client.Query().Where(
		field.Selectors[predicate.BookShelf](
			request.GetId().Selector(bookshelf.FieldID),
			request.GetName().Selector(bookshelf.FieldName),
		)...,
	)
	if request.GetWithBooks() {
		query = query.WithBooks()
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
		Data:      xslices.Trans(data, BookShelfToProto),
		Paginator: request.GetPaginator(),
	}, nil
}
