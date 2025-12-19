package bookshelfservice

import (
	"context"

	"github.com/syralon/coconut/example/internal/message"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type createController struct {
	*Dependency
}

func (c *createController) Create(ctx context.Context, request *example.CreateBookShelfRequest) (*example.CreateBookShelfResponse, error) {
	create := c.client.Create().
		SetName(request.GetName())
	if len(request.GetBookIds()) > 0 {
		create.AddRelBookIDs(request.GetBookIds()...)
	}
	data, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return &example.CreateBookShelfResponse{Data: message.BookShelfToProto(data)}, nil
}
