package service

import (
	"context"
	"errors"
	"github.com/juanfcgarcia/gostori/internal/domain"
	"github.com/juanfcgarcia/gostori/internal/logging"
	"github.com/juanfcgarcia/gostori/internal/mocks"
	"github.com/juanfcgarcia/gostori/internal/notification"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	ctrl                *gomock.Controller
	mockTransactionRepo *mocks.MockTransactionRepository
	mockAccountRepo     *mocks.MockAccountRepository
	mockNotifier        *mocks.MockNotifier
	mockFileReader      *mocks.MockReader
	fp                  *FileProcessor
)

func setup(t *testing.T) {
	ctrl = gomock.NewController(t)
	mockTransactionRepo = mocks.NewMockTransactionRepository(ctrl)
	mockAccountRepo = mocks.NewMockAccountRepository(ctrl)
	mockNotifier = mocks.NewMockNotifier(ctrl)
	mockFileReader = mocks.NewMockReader(ctrl)
	logger, _ := logging.SetupLogger()

	fp = NewFileProcessor(mockTransactionRepo, mockAccountRepo, mockNotifier, mockFileReader, 1000, logger)
}

func teardown() {
	ctrl.Finish()
}

func TestFileProcessor_ProcessFile(t *testing.T) {
	testCases := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func() {
				mockFileReader.EXPECT().ReadTransactions(gomock.Any(), "txns.csv").Return([]domain.Transaction{}, nil)
				mockAccountRepo.EXPECT().CreateAccount(gomock.Any(), "test@example.com").Return(int32(1), nil)
				mockTransactionRepo.EXPECT().GetTransactionsByAccountID(gomock.Any(), int32(1)).Return([]domain.Transaction{}, nil)
				mockNotifier.EXPECT().SendSummaryEmail(notification.SendSummaryEmailParams{
					Email:               "test@example.com",
					TotalBalance:        0,
					AverageCredit:       0,
					AverageDebit:        0,
					TransactionsByMonth: map[string]notification.MonthlySummary{},
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "file reader error",
			setupMocks: func() {
				mockFileReader.EXPECT().ReadTransactions(gomock.Any(), "txns.csv").Return(nil, errors.New("file read error"))
			},
			expectedError: errors.New("file read error"),
		},
		{
			name: "account creation error",
			setupMocks: func() {
				mockFileReader.EXPECT().ReadTransactions(gomock.Any(), "txns.csv").Return([]domain.Transaction{}, nil)
				mockAccountRepo.EXPECT().CreateAccount(gomock.Any(), "test@example.com").Return(int32(0), errors.New("account creation error"))
			},
			expectedError: errors.New("account creation error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setup(t)
			defer teardown()

			tc.setupMocks()

			err := fp.ProcessFile(context.Background(), "txns.csv", "test@example.com")
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
