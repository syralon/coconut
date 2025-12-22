package main

import (
	"context"

	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/internal/infrastructure/dependency"

	_ "github.com/syralon/coconut/toolkit/sqlite3"
)

func main() {
	client, clean, err := dependency.NewEntClient(&config.Config{Database: config.Database{Driver: "sqlite3", Source: "example.db?_pragma=foreign_keys%3Don"}})
	if err != nil {
		panic(err)
	}
	defer clean()
	if err = client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}
}
