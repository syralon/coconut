package configuration

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"strings"

	"sigs.k8s.io/yaml"
)

var (
	drivers = map[string]Driver{
		"":      LocalFileDriver,
		"local": LocalFileDriver,
	}

	unmarshalers = map[string]Unmarshaler{
		"json": json.Unmarshal,
		"yaml": func(data []byte, v any) error { return yaml.Unmarshal(data, v) },
	}
)

func SetDriver(name string, driver Driver) {
	drivers[name] = driver
}

func SetUnmarshaler(name string, unmarshaler Unmarshaler) {
	unmarshalers[name] = unmarshaler
}

func Read(ctx context.Context, v any) error {
	if !flag.Parsed() {
		flag.Parse()
	}

	driver := drivers[options.driver]
	if driver == nil {
		return fmt.Errorf("unknown config driver: %s", options.driver)
	}
	reader, err := driver.Build(ctx, options.script)
	if err != nil {
		return err
	}
	data, err := reader.Read(ctx, options.key)
	if err != nil {
		return err
	}
	ext := strings.ToLower(strings.Trim(path.Ext(options.key), "."))
	unmarshaler := unmarshalers[ext]
	if unmarshaler == nil {
		return fmt.Errorf("unknown config unmarshaler: %s", ext)
	}
	return unmarshaler(data, v)
}

func ReadT[T any](ctx context.Context) (*T, error) {
	t := new(T)
	err := Read(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, err
}
