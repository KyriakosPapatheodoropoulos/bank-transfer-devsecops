package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/KyriakosPapatheodoropoulos/bank-transfer-devsecops/app/internal/model"
	"github.com/KyriakosPapatheodoropoulos/bank-transfer-devsecops/app/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TransferHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req model.TransferRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		if req.SourceAccountID == "" ||
			req.TargetAccountID == "" ||
			req.SourceAccountID == req.TargetAccountID ||
			req.AmountCents <= 0 ||
			req.Currency == "" {

			http.Error(w, "invalid transfer request", http.StatusBadRequest)
			return
		}

		err := repository.ExecuteTransfer(
			r.Context(),
			pool,
			req.SourceAccountID,
			req.TargetAccountID,
			req.AmountCents,
			req.Currency,
		)

		if err != nil {
			if errors.Is(err, repository.ErrInsufficientFunds) {
				http.Error(w, "insufficient funds", http.StatusConflict)
				return
			}

			http.Error(w, "transfer failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{
			"status": "transfer completed",
		})
	}
}
