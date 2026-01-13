package configuration

import (
	"context"
	"fmt"
	"path"
	"strings"
)

type Option func(m *Manager)

func WithDriver(name string, driver Driver) Option {
	return func(m *Manager) {
		m.drivers[name] = driver
	}
}

func WithUnmarshaler(name string, unmarshaler Unmarshaler) Option {
	return func(m *Manager) {
		m.unmarshalers[name] = unmarshaler
	}
}

type Manager struct {
	drivers      map[string]Driver
	unmarshalers map[string]Unmarshaler

	flags
}

func NewManager(opts ...Option) *Manager {
	m := &Manager{
		drivers: map[string]Driver{
			"":      LocalFileDriver,
			"local": LocalFileDriver,
		},
		unmarshalers: map[string]Unmarshaler{
			"json": JSONUnmarshaler,
			"yaml": YAMLUnmarshaler,
		},
	}
	for _, opt := range opts {
		opt(m)
	}
	_ = m.parseFlags()
	return m
}

func (m *Manager) ReadBytes(ctx context.Context) ([]byte, error) {
	driver := m.drivers[m.driver]
	if driver == nil {
		return nil, fmt.Errorf("unknown config driver: %s", m.driver)
	}
	reader, err := driver.Build(ctx, m.script)
	if err != nil {
		return nil, err
	}
	return reader.Read(ctx, m.key)
}

func (m *Manager) Read(ctx context.Context, v any) error {
	data, err := m.ReadBytes(ctx)
	if err != nil {
		return err
	}
	ext := strings.ToLower(strings.Trim(path.Ext(m.key), "."))
	unmarshaler := m.unmarshalers[ext]
	if unmarshaler == nil {
		return fmt.Errorf("unknown config unmarshaler: %s", ext)
	}
	return unmarshaler(data, v)
}

func Read(ctx context.Context, v any, opts ...Option) error {
	return NewManager(opts...).Read(ctx, v)
}

func ReadBytes(ctx context.Context, opts ...Option) ([]byte, error) {
	return NewManager(opts...).ReadBytes(ctx)
}
