package service

import (
	"context"
	"go.uber.org/zap"
	"sync"

	"github.com/juanfcgarcia/gostori/internal/domain"
	"github.com/juanfcgarcia/gostori/internal/file"
	"github.com/juanfcgarcia/gostori/internal/notification"
	"github.com/juanfcgarcia/gostori/internal/ports"
	"github.com/juanfcgarcia/gostori/internal/transaction"
)

type FileProcessor struct {
	transactionRepo ports.TransactionRepository
	accountRepo     ports.AccountRepository
	notifier        notification.Notifier
	fileReader      file.Reader
	batchSize       int
	logger          *zap.SugaredLogger
}

func NewFileProcessor(transactionRepo ports.TransactionRepository,
	accountRepo ports.AccountRepository,
	notifier notification.Notifier,
	fileReader file.Reader,
	batchSize int,
	logger *zap.SugaredLogger) *FileProcessor {

	return &FileProcessor{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		notifier:        notifier,
		fileReader:      fileReader,
		batchSize:       batchSize,
		logger:          logger,
	}
}

func (fp *FileProcessor) ProcessFile(ctx context.Context, filePath, email string) error {
	records, err := fp.fileReader.ReadTransactions(ctx, filePath)
	if err != nil {
		fp.logger.Errorf("Error reading transactions: %v", err)
		return err
	}

	accountID, err := fp.createAccount(ctx, email)
	if err != nil {
		return err
	}

	if err := fp.processTransactions(ctx, accountID, records); err != nil {
		return err
	}

	return fp.sendSummaryEmail(ctx, accountID, email)
}

func (fp *FileProcessor) createAccount(ctx context.Context, email string) (int32, error) {
	accountID, err := fp.accountRepo.CreateAccount(ctx, email)
	if err != nil {
		fp.logger.Errorf("Failed to create account: %v", err)
		return 0, err
	}
	return accountID, nil
}

func (fp *FileProcessor) processTransactions(ctx context.Context, accountID int32, records []domain.Transaction) error {
	var wg sync.WaitGroup
	transactionCh := make(chan domain.Transaction, 100)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go fp.worker(ctx, accountID, transactionCh, &wg)
	}

	for _, record := range records {
		transactionCh <- record
	}
	close(transactionCh)

	wg.Wait()
	return nil
}

func (fp *FileProcessor) worker(ctx context.Context, accountID int32, transactionCh <-chan domain.Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	var batch []domain.Transaction
	for record := range transactionCh {
		txn := domain.Transaction{
			AccountID:     accountID,
			Amount:        record.Amount,
			Type:          record.Type,
			TransactionAt: record.TransactionAt,
		}
		batch = append(batch, txn)

		if len(batch) >= fp.batchSize {
			fp.saveBatch(ctx, batch)
			batch = []domain.Transaction{}
		}
	}
	if len(batch) > 0 {
		fp.saveBatch(ctx, batch)
	}
}

func (fp *FileProcessor) saveBatch(ctx context.Context, batch []domain.Transaction) {
	if err := fp.transactionRepo.CreateTransactions(ctx, batch); err != nil {
		fp.logger.Errorf("Failed to create transactions: %v", err)
	}
}

func (fp *FileProcessor) sendSummaryEmail(ctx context.Context, accountID int32, email string) error {
	transactions, err := fp.transactionRepo.GetTransactionsByAccountID(ctx, accountID)
	if err != nil {
		fp.logger.Errorf("Failed to get transactions: %v", err)
		return err
	}
	summary := transaction.CalculateSummary(transactions)

	notificationSummary := make(map[string]notification.MonthlySummary)
	for k, v := range summary.TransactionsByMonth {
		notificationSummary[k] = notification.MonthlySummary{TransactionCount: v.TransactionCount}
	}

	err = fp.notifier.SendSummaryEmail(notification.SendSummaryEmailParams{
		Email:               email,
		TotalBalance:        summary.TotalBalance,
		AverageCredit:       summary.AverageCredit,
		AverageDebit:        summary.AverageDebit,
		TransactionsByMonth: notificationSummary,
	})
	if err != nil {
		fp.logger.Errorf("Failed to send summary email: %v", err)
	}

	return err
}
