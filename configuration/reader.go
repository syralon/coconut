package configuration

import (
	"bytes"
	"context"
	"crypto/md5"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Content struct {
	Key      string
	Data     []byte
	Metadata map[string]string
}

type Reader interface {
	Read(ctx context.Context, key string) ([]byte, error)
	Watch(ctx context.Context, key string) (<-chan *Content, error)
}

type LocalFileReader struct {
	root     string
	interval time.Duration
	sum      []byte
}

func NewLocalFileReader(root string, watchInterval time.Duration) *LocalFileReader {
	return &LocalFileReader{root: root, interval: watchInterval}
}

func (r *LocalFileReader) Read(_ context.Context, key string) ([]byte, error) {
	data, sum, err := r.read(filepath.Join(r.root, key))
	if err != nil {
		return nil, err
	}
	r.sum = sum
	return data, nil
}

func (r *LocalFileReader) Watch(ctx context.Context, key string) (<-chan *Content, error) {
	ch := make(chan *Content)
	go r.watch(ctx, filepath.Join(r.root, key), ch)
	return ch, nil
}

func (r *LocalFileReader) watch(ctx context.Context, filename string, ch chan<- *Content) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(r.interval):
			data, sum, err := r.read(filename)
			if err != nil {
				continue // TODO handle error
			}
			if bytes.Equal(r.sum, sum) { // unchanged
				continue
			}
			r.sum = sum
			ch <- &Content{
				Key:  filename,
				Data: data,
			}
		}
	}
}

func (r *LocalFileReader) read(filename string) ([]byte, []byte, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0o600)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	re := &hashReader{Reader: file, sum: md5.New()}
	data, err := io.ReadAll(re)
	if err != nil {
		return nil, nil, err
	}
	sum := re.sum.Sum(nil)
	return data, sum, nil
}

type hashReader struct {
	io.Reader
	sum hash.Hash
}

func (h *hashReader) Read(b []byte) (int, error) {
	n, err := h.Reader.Read(b)
	h.sum.Write(b[:n])
	return n, err
}
