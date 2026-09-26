package http

import "net/http"

func NewRouter(itemsHandler *ItemsHandler, userHandler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", itemsHandler.GetItems)
	mux.HandleFunc("POST /users/{id}/withdraw", userHandler.Withdraw)
	return mux
}
