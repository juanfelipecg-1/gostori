// internal/transaction/summary_test.go

package transaction_test

import (
	"testing"
	"time"

	"github.com/juanfcgarcia/gostori/internal/domain"
	"github.com/juanfcgarcia/gostori/internal/transaction"
	"github.com/stretchr/testify/assert"
)

func TestCalculateSummary(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		transactions []domain.Transaction
		expected     transaction.Summary
	}{
		{
			name:         "no transactions",
			transactions: []domain.Transaction{},
			expected: transaction.Summary{
				TotalBalance:        0,
				TransactionsByMonth: map[string]transaction.MonthlySummary{},
				AverageCredit:       0,
				AverageDebit:        0,
			},
		},
		{
			name: "single credit transaction",
			transactions: []domain.Transaction{
				{
					Type:          transaction.CreditTransaction,
					Amount:        100,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: transaction.Summary{
				TotalBalance: 100,
				TotalCredit:  100,
				CreditCount:  1,
				TransactionsByMonth: map[string]transaction.MonthlySummary{
					"January 2024": {TransactionCount: 1},
				},
				AverageCredit: 100,
				AverageDebit:  0,
			},
		},
		{
			name: "single debit transaction",
			transactions: []domain.Transaction{
				{
					Type:          transaction.DebitTransaction,
					Amount:        50,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: transaction.Summary{
				TotalBalance: -50,
				TotalDebit:   50,
				DebitCount:   1,
				TransactionsByMonth: map[string]transaction.MonthlySummary{
					"January 2024": {TransactionCount: 1},
				},
				AverageCredit: 0,
				AverageDebit:  50,
			},
		},
		{
			name: "multiple transactions",
			transactions: []domain.Transaction{
				{
					Type:          transaction.CreditTransaction,
					Amount:        100,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Type:          transaction.DebitTransaction,
					Amount:        50,
					TransactionAt: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Type:          transaction.CreditTransaction,
					Amount:        200,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: transaction.Summary{
				TotalBalance: 250,
				TotalCredit:  300,
				CreditCount:  2,
				TotalDebit:   50,
				DebitCount:   1,
				TransactionsByMonth: map[string]transaction.MonthlySummary{
					"January 2024":  {TransactionCount: 2},
					"February 2024": {TransactionCount: 1},
				},
				AverageCredit: 150,
				AverageDebit:  50,
			},
		},
		{
			name: "transactions in different years",
			transactions: []domain.Transaction{
				{
					Type:          transaction.CreditTransaction,
					Amount:        100,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Type:          transaction.DebitTransaction,
					Amount:        50,
					TransactionAt: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: transaction.Summary{
				TotalBalance: 50,
				TotalCredit:  100,
				CreditCount:  1,
				TotalDebit:   50,
				DebitCount:   1,
				TransactionsByMonth: map[string]transaction.MonthlySummary{
					"January 2024": {TransactionCount: 1},
					"January 2025": {TransactionCount: 1},
				},
				AverageCredit: 100,
				AverageDebit:  50,
			},
		},
		{
			name: "only debit transactions",
			transactions: []domain.Transaction{
				{
					Type:          transaction.DebitTransaction,
					Amount:        50,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Type:          transaction.DebitTransaction,
					Amount:        75,
					TransactionAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: transaction.Summary{
				TotalBalance: -125,
				TotalDebit:   125,
				DebitCount:   2,
				TransactionsByMonth: map[string]transaction.MonthlySummary{
					"January 2024": {TransactionCount: 2},
				},
				AverageCredit: 0,
				AverageDebit:  62.5,
			},
		},
		{
			name: "only credit transactions",
			transactions: []domain.Transaction{
				{
					Type:          transaction.CreditTransaction,
					Amount:        50,
					TransactionAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Type:          transaction.CreditTransaction,
					Amount:        75,
					TransactionAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: transaction.Summary{
				TotalBalance: 125,
				TotalCredit:  125,
				CreditCount:  2,
				TransactionsByMonth: map[string]transaction.MonthlySummary{
					"January 2024": {TransactionCount: 2},
				},
				AverageCredit: 62.5,
				AverageDebit:  0,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := transaction.CalculateSummary(tc.transactions)
			assert.Equal(t, tc.expected, result)
		})
	}
}
