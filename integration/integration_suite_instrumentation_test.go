package integration

import (
	"encoding/csv"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
	"os"
)

func (it *IntegrationTests) createAccount() *IntegrationTests {
	email := it.argOrFail("email").(string)

	accountId, err := it.queries.CreateAccount(it.ctx, email)
	if err != nil {
		it.Fail("Error creating account: %s", err)
	}

	it.with("accountId", accountId)

	it.exec.err = err
	it.exec.result = accountId

	return it
}

func (it *IntegrationTests) AccountExists() *IntegrationTests {
	accountId := it.argOrFail("accountId").(int32)

	account, err := it.queries.GetAccountByID(it.ctx, accountId)
	if err != nil {
		return nil
	}

	if account.ID != accountId {
		it.Fail("Account ID mismatch")
	}

	it.exec.err = err
	it.exec.result = accountId

	return it
}

func (it *IntegrationTests) AccountByEmail() *IntegrationTests {
	email := it.argOrFail("email").(string)

	account, err := it.queries.GetAccountByEmail(it.ctx, email)
	it.Assert().NoError(err)
	it.Assert().Equal(email, account.Email)

	it.exec.err = err
	it.exec.result = account

	return it
}

func (it *IntegrationTests) CreateTransaction() *IntegrationTests {
	accountId := it.argOrFail("accountId").(int32)
	amount := it.argOrFail("amount").(float64)
	transactionType := it.argOrFail("transactionType").(string)
	transactionAt := it.argOrFail("transactionAt").(pgtype.Date)

	transaction := db.CreateTransactionParams{
		AccountID:     accountId,
		Amount:        amount,
		Type:          transactionType,
		TransactionAt: transactionAt,
	}

	err := it.queries.CreateTransaction(it.ctx, transaction)
	if err != nil {
		return nil
	}
	it.Assert().NoError(err)

	return it
}

func (it *IntegrationTests) CreateTransactions() *IntegrationTests {
	accountId := it.argOrFail("accountId").(int32)
	transactions := it.argOrFail("transactions").([]db.CreateTransactionParams)

	for _, transaction := range transactions {
		transaction.AccountID = accountId
		err := it.queries.CreateTransaction(it.ctx, transaction)
		it.Assert().NoError(err)
	}

	return it
}

func (it *IntegrationTests) verifyTransactionCount() *IntegrationTests {
	filePath := it.argOrFail("filePath").(string)
	expectedCount := it.countTransactionsInFile(filePath)
	it.with("expectedCount", expectedCount)
	it.verifyTransactionCountInDB()
	return it
}

func (it *IntegrationTests) countTransactionsInFile(filePath string) int {
	file, err := os.Open(filePath)
	it.Assert().NoError(err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	it.Assert().NoError(err)

	// Subtract 1 to account for the header row
	return len(records) - 1
}

func (it *IntegrationTests) verifyTransactionCountInDB() *IntegrationTests {
	expectedCount := it.argOrFail("expectedCount").(int)
	accountID := it.argOrFail("accountId").(int32)

	storedTransactions, err := it.queries.GetTransactionsByAccountID(it.ctx, accountID)
	it.Assert().NoError(err)
	it.Assert().Equal(expectedCount, len(storedTransactions), "Mismatch in transaction count between file and database")

	return it
}

func (it *IntegrationTests) readFile() *IntegrationTests {
	filePath := "../integration/txns_test.csv"
	it.with("filePath", filePath)
	return it
}

func (it *IntegrationTests) processFile() *IntegrationTests {
	email := it.argOrFail("email").(string)
	filePath := it.argOrFail("filePath").(string)

	// Process file
	err := it.appServices.FileProcessor.ProcessFile(it.ctx, filePath, email)
	it.Assert().NoError(err)

	// Retrieve the account ID using the email
	account, err := it.queries.GetAccountByEmail(it.ctx, email)
	it.Assert().NoError(err)
	it.with("accountId", account.ID)

	storedTransactions, err := it.queries.GetTransactionsByAccountID(it.ctx, account.ID)
	it.Assert().NoError(err)
	it.Assert().Greater(len(storedTransactions), 0)

	// Verify email was sent
	it.Assert().Len(it.fakeNotifierReceiver.MsgList, 1)

	return it
}

// retrieve an arg or fail if not found
func (it *IntegrationTests) argOrFail(key string) any {
	value, found := it.arg(key)
	if !found {
		it.FailNowf("get test execution argument", "a test argument was not set: %s", key)
	}

	return value
}

// arg retrieves a single value by key and existence flag
func (it *IntegrationTests) arg(key string) (any, bool) {
	value, found := it.exec.args[key]

	return value, found
}

func (it *IntegrationTests) isOk() *IntegrationTests {
	it.Assert().NoError(it.exec.err)
	it.Assert().NotNil(it.exec.result)

	return it
}

// with adds key-value pair to the execution args
func (it *IntegrationTests) with(key string, value any) *IntegrationTests {
	it.exec.args[key] = value

	return it
}
