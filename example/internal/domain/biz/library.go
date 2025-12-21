package biz

import (
	"context"

	"github.com/syralon/coconut/example/internal/domain/repository"
	"github.com/syralon/coconut/example/proto/syralon/example"
)

// LibraryBusiness
// just a business example, has no meaning.
type LibraryBusiness struct {
	rep repository.TxRepository
}

func NewLibraryBusiness(rep repository.TxRepository) *LibraryBusiness {
	return &LibraryBusiness{rep: rep}
}

func (biz *LibraryBusiness) Create(ctx context.Context, books []*example.BookCreate, shelf *example.BookShelfCreate) error {
	return biz.rep.Tx(ctx, func(ctx context.Context, txn repository.Repository) error {
		shelf.BookIds = make([]int64, 0, len(books))
		for _, book := range books {
			item, err := txn.Book().Create(ctx, book)
			if err != nil {
				return err
			}
			shelf.BookIds = append(shelf.BookIds, item.ID)
		}
		_, err := txn.BookShelf().Create(ctx, shelf)
		return err
	})
}
