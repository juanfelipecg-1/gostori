package domain

import "time"

type Transaction struct {
	ID            int32
	AccountID     int32
	Amount        float64
	Type          string
	TransactionAt time.Time
}
