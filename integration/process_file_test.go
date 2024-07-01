package integration

import (
	"github.com/brianvoe/gofakeit/v7"
)

func (it *IntegrationTests) TestProcessFileSuccess() {
	it.with("email", gofakeit.Email())
	// given a csv file
	it.readFile()
	// process file
	it.processFile()
}

func (it *IntegrationTests) TestProcessFileSuccessAndVerifyTransactionsCount() {
	it.with("email", gofakeit.Email())
	// given a csv file
	it.readFile()
	// process file
	it.processFile()
	// verify the number of transactions
	it.verifyTransactionCount()
}
