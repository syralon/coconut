package message

import (
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

func BookShelfToProto(data *ent.BookShelf) *example.BookShelf {
	shelf := &example.BookShelf{
		Id:   data.ID,
		Name: data.Name,
	}
	if data.Edges.RelBooks != nil {
		shelf.Books = make([]*example.Book, 0, len(data.Edges.RelBooks))
	}
	for _, book := range data.Edges.RelBooks {
		shelf.Books = append(shelf.Books, BookToProto(book))
	}
	return shelf
}
