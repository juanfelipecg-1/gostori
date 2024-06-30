package main

import (
	"github.com/juanfcgarcia/gostori/internal/config"
	"github.com/juanfcgarcia/gostori/internal/logging"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml", "config.local.yaml")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Setup logger
	logger, err := logging.SetupLogger()
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	logger.Infow("configuration loaded",
		"environment", cfg.EnvironmentName,
		"db_host", cfg.DBHost,
		"smtp_host", cfg.SMTPHost,
	)

	if config.IsLocal() {
		logger.Info("Running in local environment")
	} else {
		logger.Info("Running in cloud environment")
	}

}
