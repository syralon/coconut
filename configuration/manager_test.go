package configuration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	type Config struct {
		Foo string `json:"foo" yaml:"foo"`
	}
	c := new(Config)
	if err := Read(context.Background(), c); err != nil {
		t.Error(err)
	}
	assert.Equal(t, "bar", c.Foo)
}
