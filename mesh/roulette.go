package mesh

import (
	"context"
)

type Roulette interface {
	Allocate(ctx context.Context) (int, error)
	Release(ctx context.Context) error
}
