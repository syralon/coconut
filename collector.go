package coconut

import (
	"context"
	"errors"
)

type Releaser interface {
	Shutdown(ctx context.Context) error
}

type ReleaserFunc func(ctx context.Context) error

func (f ReleaserFunc) Shutdown(ctx context.Context) error {
	return f(ctx)
}

type Releasers []Releaser

func (r Releasers) Shutdown(ctx context.Context) error {
	var errs error
	for i := len(r) - 1; i >= 0; i-- {
		errs = errors.Join(errs, r[i].Shutdown(ctx))
	}
	return errs
}
