package integration

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func (it *IntegrationTests) TestCreateTransactionSuccess() {
	it.with("email", gofakeit.Email()).
		with("amount", 22.5).
		with("transactionType", "credit")

	currentDate := time.Now()
	transactionAt := pgtype.Date{
		Time:  currentDate,
		Valid: true,
	}

	it.with("transactionAt", transactionAt)
	// given an account created
	it.createAccount()
	// then transaction is created
	it.CreateTransaction()
}
