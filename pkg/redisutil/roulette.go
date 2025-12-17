package redisutil

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RouletteOption func(*Roulette)

func WithMaxID(maxID int) RouletteOption {
	return func(m *Roulette) {
		m.maxID = maxID
	}
}

func WithTTL(ttl time.Duration) RouletteOption {
	return func(m *Roulette) {
		m.ttl = ttl
	}
}

type Roulette struct {
	name   string
	client redis.UniversalClient
	maxID  int

	ttl time.Duration

	instance string
}

func NewRoulette(name string, client redis.UniversalClient, options ...RouletteOption) *Roulette {
	const defaultMaxID = 1024

	r := &Roulette{
		name:     name,
		client:   client,
		maxID:    defaultMaxID,
		ttl:      time.Minute,
		instance: uuid.New().String(),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func (a *Roulette) Allocate(ctx context.Context) (int, error) {
	for id := 0; id <= a.maxID; id++ {
		key := fmt.Sprintf("%s:workers:%d", a.name, id)

		res, err := a.client.Eval(ctx, allocateScript, []string{key}, a.instance, a.ttl).Int()
		if err != nil {
			return -1, err
		}
		if res == 1 {
			go a.heartbeat(ctx, id)
			return id, nil
		}
	}
	return -1, fmt.Errorf("no available worker ID in range [0, %d]", a.maxID)
}

func (a *Roulette) heartbeat(ctx context.Context, workerID int) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		if err := a.renew(ctx, workerID); err != nil {
			slog.ErrorContext(ctx, "lost worker id", err)
			return
		}
	}
}

func (a *Roulette) renew(ctx context.Context, workerId int) error {
	key := fmt.Sprintf("%s:workers:%d", a.name, workerId)
	res, err := a.client.Eval(ctx, renewScript, []string{key}, a.instance, a.ttl).Int()
	if err != nil {
		return err
	}
	if res == 0 {
		return fmt.Errorf("workerId %d lost ownership", workerId)
	}
	return nil
}

func (a *Roulette) Release(ctx context.Context) error {
	return nil
}
