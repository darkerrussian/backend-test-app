package app

import (
	"backend-test-app/internal/cache"
	"backend-test-app/internal/config"
	httphandler "backend-test-app/internal/http"
	"backend-test-app/internal/logic"
	"backend-test-app/internal/repository/skinport"
	"backend-test-app/pkg/httpClient"
	"backend-test-app/pkg/logger"
	"backend-test-app/pkg/signal"
	"errors"
	"net/http"
	"time"
)

func Run(cfg config.Config) {
	log := logger.New()

	skinportHTTPClient := httpClient.New(10 * time.Second)
	skinportClient := skinport.NewClient(skinportHTTPClient)
	itemsCache := cache.NewItemsCache(time.Duration(cfg.CacheTTLMins) * time.Minute)
	itemsLogic := logic.NewItemsLogic(skinportClient, itemsCache)
	itemsHandler := httphandler.NewItemsHandler(itemsLogic)

	router := httphandler.NewRouter(itemsHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	done := signal.NewSignalHandler(log, srv)

	go func() {
		log.Info("starting http server", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server error", "error", err)
			panic(err)
		}
	}()

	<-done
	log.Info("service stopped successfully")
}
