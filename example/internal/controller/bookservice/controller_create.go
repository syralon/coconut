package bookservice

import (
	"context"

	"github.com/syralon/coconut/example/internal/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type createController struct {
	*Dependency
}

func (c *createController) Create(ctx context.Context, request *example.CreateBookRequest) (*example.CreateBookResponse, error) {
	data, err := c.client.Create().SetTitle(request.GetTitle()).SetAbstract(request.GetAbstract()).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &example.CreateBookResponse{Data: message.BookToProto(data)}, nil
}
