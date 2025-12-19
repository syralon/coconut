package bookshelfservice

import (
	"github.com/syralon/coconut/example/internal/api/controller/helper"
	"github.com/syralon/coconut/example/internal/domain/repository"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type Controller struct {
	rep repository.BookShelfRepository
}

type BookShelfService struct {
	example.UnsafeBookShelfServiceServer

	createController
	listController
	updateController
}

var _ example.BookShelfServiceServer = (*BookShelfService)(nil)

func NewBookShelfService(rep repository.TxRepository) *BookShelfService {
	ctrl := &Controller{rep: rep.BookShelf()}
	s := &BookShelfService{}
	helper.SetupControllers(s, ctrl)
	return s
}
