// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID    int32
	Email string
}

type Transaction struct {
	ID            int32
	AccountID     int32
	Amount        float64
	Type          string
	TransactionAt pgtype.Date
}
