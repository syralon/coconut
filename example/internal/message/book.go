package message

import (
	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

func BookToProto(data *ent.Book) *example.Book {
	return &example.Book{
		Id:       data.ID,
		Title:    data.Title,
		Abstract: data.Abstract,
	}
}
