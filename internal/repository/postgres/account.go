package postgres

import (
	"context"
	"github.com/juanfcgarcia/gostori/internal/domain"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
)

type AccountRepository struct {
	queries *db.Queries
}

func NewAccountRepository(queries *db.Queries) *AccountRepository {
	return &AccountRepository{queries: queries}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, email string) (int32, error) {
	id, err := r.queries.CreateAccount(ctx, email)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, id int32) (*domain.Account, error) {
	account, err := r.queries.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Account{
		ID:    account.ID,
		Email: account.Email,
	}, nil
}
