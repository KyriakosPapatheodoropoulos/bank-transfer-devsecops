package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KyriakosPapatheodoropoulos/bank-transfer-devsecops/app/internal/handler"
	"github.com/KyriakosPapatheodoropoulos/bank-transfer-devsecops/app/internal/repository"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	ctx := context.Background()

	dbPool, err := repository.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	log.Println("connected to PostgreSQL")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "bank-transfer-api is running")
	})

	http.HandleFunc("/api/v1/transfer", handler.TransferHandler(dbPool))

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("server starting on :8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	signalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-signalContext.Done()

	log.Println("shutdown signal received")

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}
