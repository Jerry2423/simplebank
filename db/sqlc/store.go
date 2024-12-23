package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Queries: New(db), db: db}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// begin a seession
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// perfomr operation
	txQuery := New(tx) // build a query interface for the transaction
	err = fn(txQuery)

	// check success, not success, rollback
	if err != nil {
		rbError := tx.Rollback()
		if rbError != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbError)
		}
		return err
	}

	return tx.Commit()

}

type TransferArg struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// must be positive
	Amount int64 `json:"amount"`
}

type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Money transfer transaction in db, add transfer, add entry, update balance
func (store *Store) TransferTx(ctx context.Context, args TransferArg) (TransferResult, error) {
	var result TransferResult
	err := store.execTx(ctx, func(q *Queries) error {
		var txError error
		// create transfer
		result.Transfer, txError = q.CreateTransfer(ctx, CreateTransferParams{FromAccountID: args.FromAccountID, ToAccountID: args.ToAccountID, Amount: args.Amount})

		if txError != nil {
			return txError
		}

		// create entry 1
		result.FromEntry, txError = q.CreateEntry(ctx, CreateEntryParams{AccountID: args.FromAccountID, Amount: -args.Amount})

		if txError != nil {
			return txError
		}

		// create entry 2
		result.ToEntry, txError = q.CreateEntry(ctx, CreateEntryParams{AccountID: args.ToAccountID, Amount: args.Amount})
		if txError != nil {
			return txError
		}

		// update balance of account 1 and account 2

		// order matters, smaller id should always be first
		if args.FromAccountID < args.ToAccountID {
			// code clean up with add money
			result.FromAccount, result.ToAccount, txError = addMoney(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
			if txError != nil {
				return txError
			}
		} else {
			result.ToAccount, result.FromAccount, txError = addMoney(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
			if txError != nil {
				return txError
			}
		}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, account1 int64, amount1 int64, account2 int64, amount2 int64) (resultAcc1 Account, resultAcc2 Account, err error) {
	resultAcc1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: account1, Amount: amount1}) 
	if err != nil {
		return
	}
	resultAcc2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: account2, Amount: amount2})
	if err != nil {
		return
	}
	return
}
