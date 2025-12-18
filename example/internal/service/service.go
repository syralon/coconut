package service

import "github.com/syralon/coconut/example/ent"

type Services struct {
	Book      *BookService
	BookShelf *BookShelfService
}

func NewServices(client *ent.Client) *Services {
	return &Services{
		Book:      NewBookService(client),
		BookShelf: NewBookShelfService(client),
	}
}
