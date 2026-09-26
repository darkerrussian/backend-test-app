package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"backend-test-app/internal/entity"
)

type UserHandler struct {
	userLogic entity.UserLogic
}

func NewUserHandler(userLogic entity.UserLogic) *UserHandler {
	return &UserHandler{userLogic: userLogic}
}

type withdrawRequest struct {
	Amount string `json:"amount"`
}

type withdrawResponse struct {
	UserID     int64  `json:"user_id"`
	OldBalance string `json:"old_balance"`
	NewBalance string `json:"new_balance"`
	Amount     string `json:"amount"`
	CreatedAt  string `json:"created_at"`
}

func (h *UserHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req withdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	history, err := h.userLogic.Withdraw(r.Context(), userID, amount)
	switch {
	case err == nil:
		// продолжаем ниже
	case errors.Is(err, entity.ErrInvalidAmount):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, entity.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, entity.ErrInsufficientBalance):
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawResponse{
		UserID:     history.UserId,
		OldBalance: history.OldBalance.String(),
		NewBalance: history.NewBalance.String(),
		Amount:     history.Amount.String(),
		CreatedAt:  history.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
