package app

import (
	"backend-test-app/internal/cache"
	"backend-test-app/internal/config"
	httphandler "backend-test-app/internal/http"
	"backend-test-app/internal/logic"
	"backend-test-app/internal/repository/postgres"
	"backend-test-app/internal/repository/skinport"
	"backend-test-app/pkg/httpClient"
	"backend-test-app/pkg/logger"
	entitypostgres "backend-test-app/pkg/postgres"
	"backend-test-app/pkg/signal"
	"context"
	"errors"
	"net/http"
	"time"
)

func Run(cfg config.Config) {
	log := logger.New()

	pgPool, err := entitypostgres.NewPool(context.Background(), cfg.PostgresDSN)
	if err != nil {
		log.Error("cant connect to postgres", "error", err)
		panic(err)
	}
	defer pgPool.Close()

	skinportHTTPClient := httpClient.New(10 * time.Second)

	// repo and cache
	var (
		skinportClient = skinport.NewClient(skinportHTTPClient)
		itemsCache     = cache.NewItemsCache(time.Duration(cfg.CacheTTLMins) * time.Minute)
		userRepository = postgres.NewUserRepository(pgPool)
	)
	// logic
	var (
		itemsLogic = logic.NewItemsLogic(skinportClient, itemsCache)
		userLogic  = logic.NewUserLogic(userRepository)
	)

	// http
	var (
		itemsHandler = httphandler.NewItemsHandler(itemsLogic)
		userHandler  = httphandler.NewUserHandler(userLogic)

		router = httphandler.NewRouter(itemsHandler, userHandler)
	)

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
