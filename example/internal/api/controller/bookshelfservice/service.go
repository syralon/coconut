package bookshelfservice

import (
	"github.com/syralon/coconut/example/internal/domain/repository"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type BookShelfService struct {
	example.UnsafeBookShelfServiceServer

	rep repository.BookShelfRepository
}

var _ example.BookShelfServiceServer = (*BookShelfService)(nil)

func NewBookShelfService(rep repository.TxRepository) *BookShelfService {
	s := &BookShelfService{rep: rep.BookShelf()}
	return s
}
