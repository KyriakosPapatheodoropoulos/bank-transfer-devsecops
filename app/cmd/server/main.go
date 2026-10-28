package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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

	http.HandleFunc("/api/v1/transfer", handler.TransferHandler)

	log.Println("server starting on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
