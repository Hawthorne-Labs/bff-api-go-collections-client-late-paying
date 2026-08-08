// Package main is the entry point for bff-api-go-collections-client-late-paying.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/application/usecases"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/infrastructure"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/infrastructure/fieldcrypto"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/interface/api"
	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/interface/api/handlers"
)

func main() {
	// Load configuration
	cfg := infrastructure.LoadConfig()

	// Create core client
	coreClient, err := infrastructure.NewCoreClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create core client: %v", err)
	}

	// Create use cases
	collectionsUC := usecases.NewCollectionsUseCase(coreClient)

	// Create crypto session handler
	cryptoSessionH := handlers.NewCryptoSessionHandler(
		fieldcrypto.NewSessionManager(
			fieldcrypto.NewSessionStore(cfg.CryptoSessionTTL),
			cfg.CryptoSessionSecret, cfg.CryptoSessionIssuer, cfg.CryptoSessionTTL,
		),
	)

	// Create router
	router := api.NewRouter(collectionsUC, coreClient, cryptoSessionH)

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting bff-api-go-collections-client-late-paying on %s", addr)

	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")
}
