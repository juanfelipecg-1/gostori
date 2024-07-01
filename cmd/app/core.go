package app

import (
	"context"
	"fmt"
	"github.com/juanfcgarcia/gostori/internal/ports"

	"github.com/juanfcgarcia/gostori/internal/environment"
	"github.com/juanfcgarcia/gostori/internal/logging"
	"github.com/juanfcgarcia/gostori/internal/repository/db_factory"
	"go.uber.org/zap"
)

type CoreResources struct {
	Logger          *zap.SugaredLogger
	Config          *environment.Config
	TransactionRepo ports.TransactionRepository
	AccountRepo     ports.AccountRepository
}

func SetupCoreResources(ctx context.Context, configPath string, localConfigPath string) (*CoreResources, error) {
	// Load configuration
	cfg, err := environment.LoadConfig(configPath, localConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Setup logger
	logger, err := logging.SetupLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %v", err)
	}

	// Initialize the repository
	repo, err := db_factory.NewRepository(ctx, cfg, logger, nil)
	if err != nil {
		logger.Fatalf("failed to initialize repository: %v", err)
		return nil, err
	}

	core := &CoreResources{
		Logger:          logger,
		Config:          cfg,
		TransactionRepo: repo.TransactionRepository,
		AccountRepo:     repo.AccountRepository,
	}

	logger.Info("Core resources initialized successfully")
	return core, nil
}
