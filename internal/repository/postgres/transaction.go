package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juanfcgarcia/gostori/internal/domain"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
)

type TransactionRepository struct {
	queries *db.Queries
}

func NewTransactionRepository(queries *db.Queries) *TransactionRepository {
	return &TransactionRepository{queries: queries}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) error {
	transactionAt := pgtype.Date{
		Time:  transaction.TransactionAt,
		Valid: true,
	}

	return r.queries.CreateTransaction(ctx, db.CreateTransactionParams{
		AccountID:     transaction.AccountID,
		Amount:        transaction.Amount,
		Type:          transaction.Type,
		TransactionAt: transactionAt,
	})
}

func (r *TransactionRepository) CreateTransactions(ctx context.Context, transactions []domain.Transaction) error {
	accountIDs, amounts, types, transactionAts := initTransactionsMap(transactions)

	for i, transaction := range transactions {
		accountIDs[i] = transaction.AccountID
		amounts[i] = transaction.Amount
		types[i] = transaction.Type
		transactionAts[i] = pgtype.Date{
			Time:  transaction.TransactionAt,
			Valid: true,
		}
	}

	return r.queries.CreateTransactions(ctx, db.CreateTransactionsParams{
		Column1: accountIDs,
		Column2: amounts,
		Column3: types,
		Column4: transactionAts,
	})
}

func (r *TransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountID int32) ([]domain.Transaction, error) {
	transactions, err := r.queries.GetTransactionsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var result []domain.Transaction
	for _, t := range transactions {
		transactionAt := t.TransactionAt.Time

		result = append(result, domain.Transaction{
			ID:            t.ID,
			AccountID:     t.AccountID,
			Amount:        t.Amount,
			Type:          t.Type,
			TransactionAt: transactionAt,
		})
	}
	return result, nil
}

func initTransactionsMap(transactions []domain.Transaction) ([]int32, []float64, []string, []pgtype.Date) {
	accountIDs := make([]int32, len(transactions))
	amounts := make([]float64, len(transactions))
	types := make([]string, len(transactions))
	transactionAts := make([]pgtype.Date, len(transactions))
	return accountIDs, amounts, types, transactionAts
}
