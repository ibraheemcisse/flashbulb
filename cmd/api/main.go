package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ibraheemcisse/flashbulb/internal/config"
	"github.com/ibraheemcisse/flashbulb/internal/feed"
	"github.com/ibraheemcisse/flashbulb/internal/platform"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := platform.NewDatabase(ctx, platform.DatabaseConfig{
		URL:             cfg.DatabaseURL,
		MaxConns:        cfg.DBMaxConns,
		MinConns:        cfg.DBMinConns,
		MaxConnLifetime: cfg.DBMaxConnLifetime,
	})
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		platform.HealthHandler(db),
	)

	feedHandler := feed.NewHandler(db)

	mux.HandleFunc(
		"GET /feed",
		feedHandler.GetFeed,
	)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Flashbulb API listening on %s", cfg.HTTPAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}

	case <-signalCtx.Done():
		log.Println("Shutdown signal received")

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)

			if closeErr := server.Close(); closeErr != nil {
				log.Printf("Forced server close failed: %v", closeErr)
			}
		}

		log.Println("Flashbulb API stopped")
	}
}
