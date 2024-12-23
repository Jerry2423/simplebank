package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/Jerry2423/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func CreateRandAccount(t *testing.T) (Account, error) {
	accountParas := CreateAccountParams {Owner: util.RandOwner(), Balance: util.RandMoney(), Currency: util.RandCurrency()}
	result, err := testQueries.CreateAccount(context.Background(), accountParas)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, accountParas.Owner, result.Owner)	
	assert.Equal(t, accountParas.Balance, result.Balance)	
	assert.Equal(t, accountParas.Currency, result.Currency)	
	assert.NotZero(t, result.ID)
	assert.NotZero(t, result.CreatedAt)
	return result, err
}

func TestCreateAccount(t *testing.T) {
	CreateRandAccount(t)
}

func TestGetAccount(t *testing.T) {
	expAccount, _ := CreateRandAccount(t)
	result, err := testQueries.GetAccount(context.Background(), expAccount.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, expAccount.Owner, result.Owner)	
	assert.Equal(t, expAccount.Balance, result.Balance)	
	assert.Equal(t, expAccount.Currency, result.Currency)	
	assert.NotZero(t, result.ID)
	assert.NotZero(t, result.CreatedAt)
	assert.WithinDuration(t, expAccount.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	expAccount, _ := CreateRandAccount(t)
	accountParas := UpdateAccountParams{ID: expAccount.ID, Balance: util.RandMoney()}
	result, err := testQueries.UpdateAccount(context.Background(), accountParas)
	expAccount.Balance = accountParas.Balance

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, expAccount.Owner, result.Owner)	
	assert.Equal(t, expAccount.Balance, result.Balance)	
	assert.Equal(t, expAccount.Currency, result.Currency)	
	assert.NotZero(t, result.ID)
	assert.NotZero(t, result.CreatedAt)
	assert.WithinDuration(t, expAccount.CreatedAt, result.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	expAccount, _ := CreateRandAccount(t)
	err := testQueries.DeleteAccount(context.Background(), expAccount.ID)
	assert.NoError(t, err)
	var result Account
	result, err = testQueries.GetAccount(context.Background(), expAccount.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, result)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandAccount(t)
	}	
	funcParams := ListAccountsParams{Limit: 5, Offset: 5}
	accounts, err := testQueries.ListAccounts(context.Background(), funcParams)
	
	assert.NoError(t, err)
	fmt.Println(accounts)
	assert.Len(t, accounts, 5)
	for _, i := range accounts {
		assert.NotEmpty(t, i)
	}
}