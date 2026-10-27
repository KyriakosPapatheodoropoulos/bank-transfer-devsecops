package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KyriakosPapatheodoropoulos/bank-transfer-devsecops/app/internal/model"
)

func TransferHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "transfer request accepted",
	})
}
