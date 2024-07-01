package file_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/juanfcgarcia/gostori/internal/domain"
	"github.com/juanfcgarcia/gostori/internal/file"
	"github.com/stretchr/testify/assert"
)

func createCSVFile(t *testing.T, content string) string {
	file, err := os.CreateTemp("", "transactions_*.csv")
	assert.NoError(t, err)

	_, err = file.WriteString(content)
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	return file.Name()
}

func removeCSVFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err)
}

func TestCSVReader_ReadTransactions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		csvContent    string
		expected      []domain.Transaction
		expectedError bool
	}{
		{
			name:          "empty file",
			csvContent:    "ID,Date,Amount\n",
			expected:      []domain.Transaction{},
			expectedError: false,
		},
		{
			name:       "single credit transaction",
			csvContent: "ID,Date,Amount\n1,1/2,+100.00\n",
			expected: []domain.Transaction{
				{
					ID:            1,
					TransactionAt: time.Date(time.Now().Year(), 1, 2, 0, 0, 0, 0, time.UTC),
					Type:          "credit",
					Amount:        100.00,
				},
			},
			expectedError: false,
		},
		{
			name:       "single debit transaction",
			csvContent: "ID,Date,Amount\n2,2/3,-50.00\n",
			expected: []domain.Transaction{
				{
					ID:            2,
					TransactionAt: time.Date(time.Now().Year(), 2, 3, 0, 0, 0, 0, time.UTC),
					Type:          "debit",
					Amount:        50.00,
				},
			},
			expectedError: false,
		},
		{
			name:       "multiple transactions",
			csvContent: "ID,Date,Amount\n1,1/2,+100.00\n2,2/3,-50.00\n",
			expected: []domain.Transaction{
				{
					ID:            1,
					TransactionAt: time.Date(time.Now().Year(), 1, 2, 0, 0, 0, 0, time.UTC),
					Type:          "credit",
					Amount:        100.00,
				},
				{
					ID:            2,
					TransactionAt: time.Date(time.Now().Year(), 2, 3, 0, 0, 0, 0, time.UTC),
					Type:          "debit",
					Amount:        50.00,
				},
			},
			expectedError: false,
		},
		{
			name:          "invalid date format",
			csvContent:    "ID,Date,Amount\n1,invalid_date,+100.00\n",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid amount format",
			csvContent:    "ID,Date,Amount\n1,1/2,invalid_amount\n",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid ID format",
			csvContent:    "ID,Date,Amount\ninvalid_id,1/2,+100.00\n",
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			filePath := createCSVFile(t, tc.csvContent)
			defer removeCSVFile(t, filePath)

			reader := file.NewCSVReader()
			ctx := context.Background()

			transactions, err := reader.ReadTransactions(ctx, filePath)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, transactions)
			}
		})
	}
}
