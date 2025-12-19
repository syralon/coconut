package repository

import "context"

type Repository interface {
	Book() BookRepository
	BookShelf() BookShelfRepository
}

type Tx interface {
	Tx(ctx context.Context, fns ...func(ctx context.Context, txn Repository) error) error
}

type TxRepository interface {
	Repository
	Tx
}
