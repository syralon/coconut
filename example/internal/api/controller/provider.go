package controller

import (
	"github.com/google/wire"
	"github.com/syralon/coconut/example/internal/api/controller/bookservice"
	"github.com/syralon/coconut/example/internal/api/controller/bookshelfservice"
	"github.com/syralon/coconut/example/internal/domain/repository"
)

type Services struct {
	BookService      *bookservice.BookService
	BookShelfService *bookshelfservice.BookShelfService
}

func NewServices(rep repository.TxRepository) *Services {
	return &Services{
		BookService:      bookservice.NewBookService(rep),
		BookShelfService: bookshelfservice.NewBookShelfService(rep),
	}
}

var ProviderSet = wire.NewSet(
	NewServices,
)
