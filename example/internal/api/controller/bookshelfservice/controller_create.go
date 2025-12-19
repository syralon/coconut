package bookshelfservice

import (
	"context"

	"github.com/syralon/coconut/example/internal/api/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type createController struct {
	*Controller
}

func (c *createController) Create(ctx context.Context, request *example.CreateBookShelfRequest) (*example.CreateBookShelfResponse, error) {
	data, err := c.rep.Create(ctx, request.GetCreate())
	if err != nil {
		return nil, err
	}
	return &example.CreateBookShelfResponse{Data: message.BookShelfToProto(data)}, nil
}
