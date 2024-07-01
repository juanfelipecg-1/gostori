package ports

import (
	"context"
	"github.com/juanfcgarcia/gostori/internal/domain"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction domain.Transaction) error
	CreateTransactions(ctx context.Context, transaction []domain.Transaction) error
	GetTransactionsByAccountID(ctx context.Context, accountID int32) ([]domain.Transaction, error)
}
