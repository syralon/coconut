package bookshelfservice

import (
	"context"

	"github.com/syralon/coconut/example/proto/syralon/example"
)

type updateController struct {
	*Controller
}

func (c *updateController) Update(ctx context.Context, request *example.UpdateBookShelfRequest) (*example.UpdateBookShelfResponse, error) {
	if err := c.rep.Update(ctx, request.GetId(), request.GetUpdate()); err != nil {
		return nil, err
	}
	return &example.UpdateBookShelfResponse{}, nil
}
