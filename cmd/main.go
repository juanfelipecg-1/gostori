package main

import (
	"context"
	"github.com/juanfcgarcia/gostori/cmd/app"
	"github.com/juanfcgarcia/gostori/cmd/flags"
	"github.com/juanfcgarcia/gostori/cmd/services"
	"log"
)

func main() {
	filePath, email := flags.ParseFlags()

	// Create context
	ctx := context.Background()

	// Initialize core resources
	core, err := app.SetupCoreResources(ctx, "config.yaml", "config.local.yaml")
	if err != nil {
		log.Fatalf("failed to setup core resources: %v", err)
	}

	// Initialize application services
	appServices, err := services.SetupAppServices(core)
	if err != nil {
		core.Logger.Fatalf("failed to setup application services: %v", err)
	}

	// Process file
	err = appServices.FileProcessor.ProcessFile(ctx, filePath, email)
	if err != nil {
		core.Logger.Fatalf("failed to process transactions: %v", err)
	}

	core.Logger.Info("Transactions processed successfully")
}
