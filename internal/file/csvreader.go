package file

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/juanfcgarcia/gostori/internal/domain"
)

type CSVReader struct{}

func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

func (r *CSVReader) ReadTransactions(ctx context.Context, filePath string) ([]domain.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, 0)
	currentYear := time.Now().Year()

	for _, record := range records[1:] { // Remove header
		id, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(record[2][1:], 64)
		if err != nil {
			return nil, err
		}
		date, err := time.Parse("1/2", record[1])
		if err != nil {
			return nil, err
		}

		date = date.AddDate(currentYear-date.Year(), 0, 0)

		transactionType := "debit"
		if record[2][:1] == "+" {
			transactionType = "credit"
		}

		transaction := domain.Transaction{
			ID:            int32(id),
			TransactionAt: date,
			Type:          transactionType,
			Amount:        amount,
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
