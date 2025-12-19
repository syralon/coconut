package bookshelfservice

import (
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/internal/controller/helper"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

type Dependency struct {
	client *ent.BookShelfClient
}

type BookShelfService struct {
	example.UnsafeBookShelfServiceServer

	createController
	listController
}

var _ example.UnsafeBookShelfServiceServer = (*BookShelfService)(nil)

func NewBookShelfService(client *ent.Client) *BookShelfService {
	dep := &Dependency{client: client.BookShelf}
	s := &BookShelfService{}
	helper.SetupController(s, dep)
	return s
}
