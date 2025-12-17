package xio

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitWriter(t *testing.T) {
	{
		buf := bytes.NewBuffer(nil)
		w := NewLimitWriter(buf, 5)
		n, err := w.Write([]byte("test"))
		assert.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, []byte("test"), buf.Bytes())
	}
	{
		buf := bytes.NewBuffer(nil)
		w := NewLimitWriter(buf, 5)
		n, err := w.Write([]byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte("hello"), buf.Bytes())
	}

	{
		buf := bytes.NewBuffer(nil)
		w := NewLimitWriter(buf, 5)
		n, err := w.Write([]byte("hello world"))
		assert.Equal(t, true, errors.Is(err, ErrUpToLimit))
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte("hello"), buf.Bytes())
	}
}
