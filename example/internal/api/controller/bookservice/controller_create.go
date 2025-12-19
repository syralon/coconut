package bookservice

import (
	"context"

	"github.com/syralon/coconut/example/internal/api/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type createController struct {
	*Controller
}

func (c *createController) Create(ctx context.Context, request *example.CreateBookRequest) (*example.CreateBookResponse, error) {
	data, err := c.rep.Create(ctx, request.GetCreate())
	if err != nil {
		return nil, err
	}
	return &example.CreateBookResponse{Data: message.BookToProto(data)}, nil
}
