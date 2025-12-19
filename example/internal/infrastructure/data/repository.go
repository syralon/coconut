package data

import (
	"context"
	stderrors "errors"

	"github.com/syralon/coconut/example/ent"
	"github.com/syralon/coconut/example/internal/domain/repository"
	"github.com/syralon/coconut/proto/syralon/coconut/errors"
)

type Repository struct {
	c         *ent.Client
	book      *BookRepository
	bookShelf *BookShelfRepository
}

func (rep *Repository) Book() repository.BookRepository {
	if rep.book == nil {
		rep.book = NewBookRepository(rep.c)
	}
	return rep.book
}

func (rep *Repository) BookShelf() repository.BookShelfRepository {
	if rep.bookShelf == nil {
		rep.bookShelf = NewBookShelfRepository(rep.c)
	}
	return rep.bookShelf
}

func (rep *Repository) Tx(ctx context.Context, fns ...func(ctx context.Context, txn repository.Repository) error) (err error) {
	var tx *ent.Tx
	tx, err = rep.c.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rec := errors.Recovery(recover()); rec != nil {
			err = rec
		}
		if err == nil {
			err = tx.Commit()
		} else {
			err = stderrors.Join(err, tx.Rollback())
		}
	}()

	txn := &Repository{c: tx.Client()}
	for _, fn := range fns {
		if err = fn(ctx, txn); err != nil {
			return err
		}
	}
	return nil
}

func NewRepository(client *ent.Client) repository.TxRepository {
	return &Repository{
		c:         client,
		book:      NewBookRepository(client),
		bookShelf: NewBookShelfRepository(client),
	}
}
