package xio

import (
	"bytes"
	"errors"
	"io"
)

var ErrUpToLimit = errors.New("write bytes up to limit")

type LimitWriter struct {
	w       io.Writer
	limit   int
	written int
}

func NewLimitWriter(w io.Writer, limit int) *LimitWriter {
	return &LimitWriter{
		w:       w,
		limit:   limit,
		written: 0,
	}
}

func (w *LimitWriter) Write(p []byte) (n int, err error) {
	left := w.limit - w.written
	if left <= 0 {
		return 0, ErrUpToLimit
	}
	write := len(p)
	if write > left {
		n, err = w.w.Write(p[:left])
		if err == nil {
			err = ErrUpToLimit
		} else {
			err = errors.Join(err, ErrUpToLimit)
		}
	} else {
		n, err = w.w.Write(p[:write])
	}
	w.written += n
	return n, err
}

type LimitBuffer struct {
	buf *bytes.Buffer
	*LimitWriter
}

func NewLimitBuffer(limit int) *LimitBuffer {
	buf := bytes.NewBuffer(nil)
	return &LimitBuffer{
		buf:         buf,
		LimitWriter: NewLimitWriter(buf, limit),
	}
}

func (b *LimitBuffer) Bytes() []byte {
	return b.buf.Bytes()
}
