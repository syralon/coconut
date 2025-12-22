package bookservice

import (
	"context"

	"github.com/syralon/coconut/example/internal/api/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/toolkit/xslices"
)

func (s *BookService) List(ctx context.Context, request *example.ListBookRequest) (*example.ListBookResponse, error) {
	data, paginator, err := s.rep.List(ctx, request.GetOptions(), request.GetPaginator())
	if err != nil {
		return nil, err
	}
	return &example.ListBookResponse{
		Data:      xslices.Trans(data, message.BookToProto),
		Paginator: paginator,
	}, nil
}
