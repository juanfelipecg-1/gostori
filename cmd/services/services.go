package services

import (
	"github.com/juanfcgarcia/gostori/cmd/app"
	"github.com/juanfcgarcia/gostori/internal/file"
	"github.com/juanfcgarcia/gostori/internal/mailer"
	"github.com/juanfcgarcia/gostori/internal/notification"
	"github.com/juanfcgarcia/gostori/internal/service"
)

type AppServices struct {
	FileProcessor *service.FileProcessor
	Notifier      notification.Notifier
	Mailer        mailer.Mailer
}

func SetupAppServices(core *app.CoreResources) (*AppServices, error) {
	// Initialize mailer
	mailSvc := mailer.NewSMTPMailer(core.Config.SMTPHost, core.Config.SMTPPort)
	notifier := notification.NewNotifier(mailSvc)

	// Initialize file reader
	fileReader := file.NewCSVReader()

	// Initialize file processor
	fileProcessor := service.NewFileProcessor(core.TransactionRepo, core.AccountRepo, notifier, fileReader, core.Config.MaxBatchSize, core.Logger)

	appServices := &AppServices{
		FileProcessor: fileProcessor,
		Notifier:      notifier,
		Mailer:        mailSvc,
	}

	core.Logger.Info("Application services initialized successfully")
	return appServices, nil
}
