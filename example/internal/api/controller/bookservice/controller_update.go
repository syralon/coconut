package bookservice

import (
	"context"

	"github.com/syralon/coconut/example/proto/syralon/example"
)

func (s *BookService) Update(ctx context.Context, request *example.UpdateBookRequest) (*example.UpdateBookResponse, error) {
	if err := s.rep.Update(ctx, request.GetId(), request.GetUpdate()); err != nil {
		return nil, err
	}
	return &example.UpdateBookResponse{}, nil
}
