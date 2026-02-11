package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"TestProject_itk/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WalletService interface {
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
	ApplyOperation(ctx context.Context, id uuid.UUID, op domain.OperationType, amount int64) (int64, error)
}

type Handler struct {
	svc WalletService
}

func NewHandler(svc WalletService) *Handler {
	return &Handler{svc: svc}
}

type postWalletRequest struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

type walletResponse struct {
	WalletID string `json:"walletId"`
	Balance  int64  `json:"balance"`
}

func (h *Handler) PostWalletOperation(w http.ResponseWriter, r *http.Request) {
	r, cancel := withTimeout(r, 3*time.Second)
	defer cancel()

	var req postWalletRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	id, err := uuid.Parse(req.WalletID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid walletId"})
		return
	}

	op := domain.OperationType(req.OperationType)
	if !op.Valid() {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid operationType"})
		return
	}

	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "amount must be > 0"})
		return
	}

	balance, err := h.svc.ApplyOperation(r.Context(), id, op, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWalletNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet not found"})
		case errors.Is(err, domain.ErrInsufficientFund):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "insufficient funds"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal"})
		}
		return
	}

	writeJSON(w, http.StatusOK, walletResponse{WalletID: id.String(), Balance: balance})
}

func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	r, cancel := withTimeout(r, 2*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "walletId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid walletId"})
		return
	}

	balance, err := h.svc.GetBalance(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWalletNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal"})
		}
		return
	}

	writeJSON(w, http.StatusOK, walletResponse{WalletID: id.String(), Balance: balance})
}
