package controller

import (
	"github.com/google/wire"
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/internal/controller/bookservice"
	"github.com/syralon/coconut/example/internal/controller/bookshelfservice"
)

type Services struct {
	BookService      *bookservice.BookService
	BookShelfService *bookshelfservice.BookShelfService
}

func NewServices(client *ent.Client) *Services {
	return &Services{
		BookService:      bookservice.NewBookService(client),
		BookShelfService: bookshelfservice.NewBookShelfService(client),
	}
}

var ProviderSet = wire.NewSet(
	NewServices,
)
