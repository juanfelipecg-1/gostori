package integration

import "github.com/brianvoe/gofakeit/v7"

func (it *IntegrationTests) TestCreateAccountSuccess() {
	it.with("email", gofakeit.Email())
	// when requesting to create an account
	it.createAccount()
	// then success
	it.AccountExists()
}

func (it *IntegrationTests) TestGetAccountByEmailSuccess() {
	it.with("email", gofakeit.Email())
	// when requesting to create an account
	it.createAccount()
	// then retrieve the account by email
	it.AccountByEmail()
}

func (it *IntegrationTests) TestCreateMultipleAccountsWithSameEmail() {
	email := gofakeit.Email()
	it.with("email", email)
	// create the first account
	it.createAccount().isOk()
	it.AccountExists()
	account1ID := it.argOrFail("accountId").(int32)

	// then a second account with the same email is created
	it.createAccount().isOk()
	it.AccountExists()
	account2ID := it.argOrFail("accountId").(int32)

	// then verify that the second account ID is different from the first
	it.Assert().NotEqual(account1ID, account2ID, "Expected different account IDs for the same email")
}
