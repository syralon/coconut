package sqlite3

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func init() {
	const src = "sqlite"
	const dst = "sqlite3"

	drivers := sql.Drivers()
	found := false
	for _, d := range drivers {
		if d == src {
			found = true
			break
		}
	}
	if !found {
		panic("sqlite driver not registered")
	}

	db, err := sql.Open(src, "")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql.Register(dst, db.Driver())
}
