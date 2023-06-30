package db

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	TransferTx(cxt context.Context, arg TransferTxParams) (TransferTxResult, error)
 	CreateUserTx(cxt context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) 
 	VerifyEmailTx(cxt context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) 
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewSQLStore ...
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (store *SQLStore) execTx(cxt context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(cxt, nil)
	if err != nil {
		return err
	}
	// Pass the transaction object to the queries object
	q := New(tx)
	err = fn(q)
	if err != nil {
		// Rollback the transaction if there is an error
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	// Commit the transaction if there is no error
	return tx.Commit()
}