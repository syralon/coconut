package bookshelfservice

import (
	"context"

	"github.com/syralon/coconut/example/internal/api/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/toolkit/xslices"
)

func (s *BookShelfService) List(ctx context.Context, request *example.ListBookShelfRequest) (*example.ListBookShelfResponse, error) {
	data, paginator, err := s.rep.List(ctx, request.GetOptions(), request.GetPaginator())
	if err != nil {
		return nil, err
	}
	return &example.ListBookShelfResponse{
		Data:      xslices.Trans(data, message.BookShelfToProto),
		Paginator: paginator,
	}, nil
}
