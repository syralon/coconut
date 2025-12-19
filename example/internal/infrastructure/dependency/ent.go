package dependency

import (
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/syralon/coconut/toolkit/sqlite3"
)

func NewEntClient(c *config.Config) (*ent.Client, func(), error) {
	client, err := ent.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	return client, func() { _ = client.Close() }, err
}
