package bookservice

import (
	"github.com/syralon/coconut/example/internal/api/controller/helper"
	"github.com/syralon/coconut/example/internal/domain/repository"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type Controller struct {
	//client *ent.BookClient

	rep repository.BookRepository
}

type BookService struct {
	example.UnsafeBookServiceServer

	createController
	listController
	updateController
}

var _ example.BookServiceServer = (*BookService)(nil)

func NewBookService(rep repository.TxRepository) *BookService {
	s := &BookService{}
	ctrl := &Controller{rep: rep.Book()}
	helper.SetupControllers(s, ctrl)
	return s
}
