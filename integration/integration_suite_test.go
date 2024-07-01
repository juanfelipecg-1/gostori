package integration

import (
	"context"
	"github.com/juanfcgarcia/gostori/cmd/services"
	"github.com/juanfcgarcia/gostori/internal/file"
	"github.com/juanfcgarcia/gostori/internal/logging"
	"github.com/juanfcgarcia/gostori/internal/mailer"
	"github.com/juanfcgarcia/gostori/internal/notification"
	"github.com/juanfcgarcia/gostori/internal/repository/db"
	"github.com/juanfcgarcia/gostori/internal/repository/db_factory"
	"github.com/juanfcgarcia/gostori/internal/service"
	"github.com/juanfcgarcia/gostori/internal/testutils/fakes"
	"path/filepath"
	"testing"

	"github.com/juanfcgarcia/gostori/internal/environment"
	"github.com/juanfcgarcia/gostori/internal/testutils"
	"github.com/stretchr/testify/suite"
)

type IntegrationTests struct {
	*suite.Suite
	tdb                  *testutils.TestingDB
	ctx                  context.Context
	queries              *db.Queries
	fakeNotifierReceiver *fakes.DummyNotifierReceiver
	notifier             notification.Notifier
	config               *environment.Config
	appServices          services.AppServices

	cleanup testutils.CloseFunc

	// fields mutable per test
	exec *testExec
}

type testExec struct {
	result any
	err    error
	args   map[string]any
}

func (it *IntegrationTests) SetupSuite() {
	it.ctx = context.Background()

	absMigratePath, err := filepath.Abs("../db-scripts/migrations")
	if err != nil {
		it.FailNow("Failed to get absolute path for migrations", err.Error())
	}

	database, cleanup, err := testutils.SetupTestDb(absMigratePath)
	if err != nil {
		it.FailNow("Failed to setup test DB", err.Error())
	}
	it.tdb = database
	it.cleanup = cleanup

	cfg, err := environment.LoadConfig("../integration/test.config.yaml")
	it.Assert().NoError(err)

	it.config = cfg
	it.queries = db.New(it.tdb.Conn)

	mailSvc := mailer.NewSMTPMailer(it.config.SMTPHost, it.config.SMTPPort)
	it.notifier = notification.NewNotifier(mailSvc)

	notifier := &fakes.DummyNotifier{
		Receiver: fakes.DummyNotifierReceiver{},
	}

	it.fakeNotifierReceiver = &notifier.Receiver

	// Setup logger
	logger, err := logging.SetupLogger()
	if err != nil {
		it.FailNow("Failed to setup logger", err.Error())
	}

	repo, err := db_factory.NewRepository(it.ctx, it.config, logger, it.tdb.Conn)
	it.Assert().NoError(err)

	fileReader := file.NewCSVReader()
	fileProcessor := service.NewFileProcessor(repo.TransactionRepository, repo.AccountRepository, notifier, fileReader, it.config.MaxBatchSize, logger)

	it.appServices = services.AppServices{
		FileProcessor: fileProcessor,
		Notifier:      notifier,
		Mailer:        mailSvc,
	}

}

func (it *IntegrationTests) TearDownSuite() {
	err := it.cleanup()
	if err != nil {
		return
	}
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, &IntegrationTests{
		Suite: new(suite.Suite),
	})
}

func (it *IntegrationTests) SetupTest() {
	// reset exec
	it.exec = &testExec{
		args: make(map[string]any),
	}

	defer func() {
		it.fakeNotifierReceiver.MsgList = make([]string, 0)
	}()
}
