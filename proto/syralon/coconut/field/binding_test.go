package field

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestBinder(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("X-Client-Id", "111", "X-Ids", "1", "X-Ids", "2", "X-Ids", "1"))
	{
		e := &Example{}
		err := Bind(ctx, e)
		assert.NoError(t, err)
		assert.Equal(t, "111", e.GetClientId())
		assert.Equal(t, []int32{1, 2, 1}, e.GetIds())
	}

	{
		b := NewBinder(WithCache(true))
		e := &Example{}
		err := b.Bind(ctx, e)
		assert.NoError(t, err)
		assert.Equal(t, "111", e.GetClientId())
		assert.Equal(t, []int32{1, 2, 1}, e.GetIds())
	}
}

func BenchmarkBinder(b *testing.B) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("X-Client-Id", "111", "X-Ids", "1", "X-Ids", "2", "X-Ids", "1"))
	e := &Example{}

	for i := 0; i < b.N; i++ {
		_ = Bind(ctx, e)
	}
}

func BenchmarkBinderWithCache(b *testing.B) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("X-Client-Id", "111", "X-Ids", "1", "X-Ids", "2", "X-Ids", "1"))
	e := &Example{}
	binder := NewBinder(WithCache(true))

	for i := 0; i < b.N; i++ {
		_ = binder.Bind(ctx, e)
	}
}
