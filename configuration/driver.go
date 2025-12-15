package configuration

import (
	"context"
	"time"
)

type Driver interface {
	Build(ctx context.Context, script string) (Reader, error)
}

var LocalFileDriver = &localFileDriver{}

type localFileDriver struct{}

func (localFileDriver) Build(_ context.Context, _ string) (Reader, error) {
	return NewLocalFileReader(".", time.Minute), nil
}
