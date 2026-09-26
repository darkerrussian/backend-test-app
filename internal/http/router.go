package http

import "net/http"

func NewRouter(itemsHandler *ItemsHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", itemsHandler.GetItems)
	return mux
}
