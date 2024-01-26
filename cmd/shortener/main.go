package main

import (
	"context"
	"github.com/pingachguk/ya-shortener/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingachguk/ya-shortener/config"
	"github.com/pingachguk/ya-shortener/internal/compresser"
	"github.com/pingachguk/ya-shortener/internal/handlers"
	"github.com/pingachguk/ya-shortener/internal/logger"
	"github.com/rs/zerolog/log"
)

func GetRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(
		compresser.CompressMiddleware,
		logger.LogMiddleware,
	)

	router.Route("/api", func(r chi.Router) {
		r.Post("/shorten", handlers.APICreateShortHandler)
	})

	router.Get("/{short}", handlers.TryRedirectHandler)
	router.Post("/", handlers.CreateShortHandler)
	router.Get("/ping", handlers.PingDatabase)

	return router
}

func closeStorage() {
	err := storage.GetStorage().Close(context.Background())
	if err != nil {
		log.Err(err).Msgf("")
	}
}

func main() {
	config.InitConfig()
	defer closeStorage()

	log.Info().Msgf("[*] Application address: %s", config.Config.Base)
	log.Info().Msgf("[*] Base address: %s", config.Config.Base)

	if err := http.ListenAndServe(config.Config.App, GetRouter()); err != nil {
		panic(err)
	}
}
