package db

import (
	"context"
	"database/sql"
)

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore ...
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (store *Store) execTx(cxt context.Context, fn func(*Queries) error) error {
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

// TransferTxParams
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, and account entries, and update accounts' balance within a single database transaction.

func (store *Store) TransferTx(cxt context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(cxt, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(cxt, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(cxt, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(cxt, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// get account -> update account balance
		result.FromAccount, err = q.AddAccountBalance(cxt, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(cxt, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
