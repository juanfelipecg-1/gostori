package file

import (
	"context"
	"github.com/juanfcgarcia/gostori/internal/domain"
)

type Reader interface {
	ReadTransactions(ctx context.Context, filePath string) ([]domain.Transaction, error)
}
