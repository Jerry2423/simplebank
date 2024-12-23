package db

import (
	"context"
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Only test deadlock scenerio, do not care about the result
func TestTransferTxDeadlock(t *testing.T) {
	testStore := NewStore(globalDB)
	// test concurrent transaction
	account1, _ := CreateRandAccount(t)
	account2, _ := CreateRandAccount(t)

	errors := make(chan error)
	var n int = 10
	var amount int64 = 10
	for i := 0; i < n; i++ {
		fromAccount := account1
		toAccount := account2
		if i%2 == 1 {
			fromAccount, toAccount = toAccount, fromAccount
		}
		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferArg{FromAccountID: fromAccount.ID, ToAccountID: toAccount.ID, Amount: amount})
			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		curr_err := <-errors
		assert.NoError(t, curr_err)
	}

	// check the final balance
	finalAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1.Balance, finalAccount1.Balance)

	var finalAccount2 Account
	finalAccount2, err = testStore.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)
	assert.Equal(t, finalAccount2.Balance, account2.Balance)

}

func TestTransferTx(t *testing.T) {
	testStore := NewStore(globalDB)
	// test concurrent transaction
	account1, _ := CreateRandAccount(t)
	account2, _ := CreateRandAccount(t)

	results := make(chan TransferResult)
	errors := make(chan error)
	var n int = 5
	var amount int64 = 10
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferArg{FromAccountID: account1.ID, ToAccountID: account2.ID, Amount: amount})
			errors <- err
			results <- result
		}()
	}

	record_k := make(map[int64]bool)
	for i := 0; i < n; i++ {
		curr_err := <-errors
		curr_result := <-results
		assert.NoError(t, curr_err)
		assert.NotEmpty(t, curr_result)

		transfer := curr_result.Transfer
		assert.NotEmpty(t, transfer)
		assert.Equal(t, transfer.Amount, amount)
		assert.Equal(t, transfer.FromAccountID, account1.ID)
		assert.Equal(t, transfer.ToAccountID, account2.ID)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)
		_, err := testStore.GetTransfer(context.Background(), transfer.ID) // check inserted success
		assert.NoError(t, err)

		// check entry
		fromEntry := curr_result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, fromEntry.Amount, -amount)
		assert.Equal(t, fromEntry.AccountID, account1.ID)
		assert.NotZero(t, fromEntry.CreatedAt)
		assert.NotZero(t, fromEntry.ID)
		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		toEntry := curr_result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, toEntry.Amount, amount)
		assert.Equal(t, toEntry.AccountID, account2.ID)
		assert.NotZero(t, toEntry.CreatedAt)
		assert.NotZero(t, toEntry.ID)
		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		// check account
		fromAccount := curr_result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, fromAccount.ID, account1.ID)

		toAccount := curr_result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, toAccount.ID, account2.ID)

		//  check balance of users
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		// same, multiple, unique
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1%amount == 0)
		k := diff1 / amount
		assert.True(t, k >= 1 && k <= int64(n))
		assert.NotContains(t, record_k, k)
		record_k[k] = true
	}

	// check the final balance
	finalAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	diff := account1.Balance - finalAccount1.Balance
	assert.Equal(t, diff, int64(n)*amount)

	var finalAccount2 Account
	finalAccount2, err = testStore.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)
	diff = finalAccount2.Balance - account2.Balance
	assert.Equal(t, diff, int64(n)*amount)
}
