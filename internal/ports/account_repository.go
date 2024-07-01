package ports

import (
	"context"
	"github.com/juanfcgarcia/gostori/internal/domain"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, email string) (int32, error)
	GetAccountByID(ctx context.Context, id int32) (*domain.Account, error)
}
