package bookservice

import (
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/ent/book"
	"github.com/syralon/coconut/example/internal/controller/helper"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

var BookOrderBy = map[example.BookOrder]string{
	example.BookOrder_BookOrderByID: book.FieldID,
}

type Dependency struct {
	client *ent.BookClient
}

type BookService struct {
	example.UnsafeBookServiceServer

	createController
	listController
}

var _ example.BookServiceServer = (*BookService)(nil)

func NewBookService(client *ent.Client) *BookService {
	dep := &Dependency{client: client.Book}
	s := &BookService{}

	helper.SetupController(s, dep)
	return s
}
