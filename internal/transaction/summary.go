// internal/transaction/summary.go

package transaction

import (
	"github.com/juanfcgarcia/gostori/internal/domain"
	"strconv"
	"time"
)

type MonthlySummary struct {
	TransactionCount int
}

type Summary struct {
	TotalBalance        float64
	TotalCredit         float64
	TotalDebit          float64
	CreditCount         int
	DebitCount          int
	TransactionsByMonth map[string]MonthlySummary
	AverageCredit       float64
	AverageDebit        float64
}

const (
	CreditTransaction = "credit"
	DebitTransaction  = "debit"
)

func CalculateSummary(transactions []domain.Transaction) Summary {
	summary := processTransactions(transactions)

	summary.AverageCredit = calculateAverage(summary.TotalCredit, summary.CreditCount)
	summary.AverageDebit = calculateAverage(summary.TotalDebit, summary.DebitCount)

	return summary
}

func processTransactions(transactions []domain.Transaction) Summary {
	summary := Summary{
		TransactionsByMonth: make(map[string]MonthlySummary),
	}

	for _, transaction := range transactions {
		monthYear := formatMonthYear(transaction.TransactionAt)
		monthlySummary := summary.TransactionsByMonth[monthYear]
		monthlySummary.TransactionCount++
		summary.TransactionsByMonth[monthYear] = monthlySummary

		switch transaction.Type {
		case CreditTransaction:
			addCreditTransaction(&summary, transaction)
		case DebitTransaction:
			addDebitTransaction(&summary, transaction)
		}
	}

	return summary
}

func addCreditTransaction(summary *Summary, transaction domain.Transaction) {
	summary.TotalBalance += transaction.Amount
	summary.TotalCredit += transaction.Amount
	summary.CreditCount++
}

func addDebitTransaction(summary *Summary, transaction domain.Transaction) {
	summary.TotalBalance -= transaction.Amount
	summary.TotalDebit += transaction.Amount
	summary.DebitCount++
}

func calculateAverage(total float64, count int) float64 {
	if count > 0 {
		return total / float64(count)
	}
	return 0.0
}

func formatMonthYear(date time.Time) string {
	month := date.Format("January")
	year := date.Year()
	return month + " " + strconv.Itoa(year)
}
