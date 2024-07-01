package integration

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
	"time"
)

func (it *IntegrationTests) TestCreateTransactionsSuccess() {
	it.with("email", gofakeit.Email())

	currentDate := time.Now()
	transactionAt := pgtype.Date{
		Time:  currentDate,
		Valid: true,
	}

	transactions := []db.CreateTransactionParams{
		{
			Amount:        22.5,
			Type:          "credit",
			TransactionAt: transactionAt,
		},
		{
			Amount:        15.2,
			Type:          "debit",
			TransactionAt: transactionAt,
		},
		{
			Amount:        3.5,
			Type:          "debit",
			TransactionAt: transactionAt,
		},
	}

	it.with("transactions", transactions)
	// given an account created
	it.createAccount()
	// then transactions are created
	it.CreateTransactions()
}
