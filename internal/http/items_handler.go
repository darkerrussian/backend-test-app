package http

import (
	"backend-test-app/internal/entity"
	"encoding/json"
	"net/http"
	"strconv"
)

type ItemsHandler struct {
	itemsLogic entity.ItemsLogic
}

func NewItemsHandler(service entity.ItemsLogic) *ItemsHandler {
	return &ItemsHandler{itemsLogic: service}
}

func (h *ItemsHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	// 730 - Counter Strike 2 skinport default
	appId := 730
	if v := r.URL.Query().Get("app_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "invalid app_id", http.StatusBadRequest)
			return
		}
		appId = parsed
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "EUR"
	}

	if !entity.IsSupportedCurrency(currency) {
		http.Error(w, "unsupported currency", http.StatusBadRequest)
		return
	}

	items, err := h.itemsLogic.GetItems(r.Context(), entity.SkinportParams{
		AppId:    appId,
		Currency: currency,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
