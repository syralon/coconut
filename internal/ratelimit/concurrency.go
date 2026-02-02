package ratelimit

import "context"

type ConcurrencyLimiter struct {
	limit chan struct{}
}

func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{make(chan struct{}, limit)}
}

func (cl *ConcurrencyLimiter) Acquire(ctx context.Context) error {
	select {
	case cl.limit <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (cl *ConcurrencyLimiter) Release() {
	<-cl.limit
}

func (cl *ConcurrencyLimiter) Do(ctx context.Context, fn func() error) error {
	if err := cl.Acquire(ctx); err != nil {
		return err
	}
	defer cl.Release()
	return fn()
}
