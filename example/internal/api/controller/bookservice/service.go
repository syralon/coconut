package bookservice

import (
	"github.com/syralon/coconut/example/internal/domain/repository"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type BookService struct {
	example.UnsafeBookServiceServer

	rep repository.BookRepository
}

var _ example.BookServiceServer = (*BookService)(nil)

func NewBookService(rep repository.TxRepository) *BookService {
	s := &BookService{
		rep: rep.Book(),
	}
	return s
}
