package bookservice

import (
	"context"

	"github.com/syralon/coconut/example/proto/syralon/example"
)

type updateController struct {
	*Controller
}

func (c *updateController) Update(ctx context.Context, request *example.UpdateBookRequest) (*example.UpdateBookResponse, error) {
	if err := c.rep.Update(ctx, request.GetId(), request.GetUpdate()); err != nil {
		return nil, err
	}
	return &example.UpdateBookResponse{}, nil
}
