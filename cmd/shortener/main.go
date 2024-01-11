package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingachguk/ya-shortener/config"
	"github.com/pingachguk/ya-shortener/internal/compresser"
	"github.com/pingachguk/ya-shortener/internal/logger"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/pingachguk/ya-shortener/internal/storage"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
)

func tryRedirectHandler(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	shorten := storage.GetStorage().GetByShort(short)

	if shorten != nil {
		http.Redirect(w, r, shorten.OriginalUrl, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func createShortHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	} else if len(body) == 0 {
		http.Error(w, "Bad request data: empty body", http.StatusBadRequest)
		return
	}

	url := string(body)
	short, err := shortid.GetDefault().Generate()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	s := models.NewShorten(short, url)
	err = storage.GetStorage().AddShorten(*s)
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", config.Config.Base, short)))
}

func apiCreateShortHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Request

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		errorResponse(w, "Bad Content-Type: need application/json", http.StatusNotAcceptable)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err == io.EOF {
		errorResponse(w, "Bad request data: empty body", http.StatusBadRequest)
		return
	} else if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	short, err := shortid.GetDefault().Generate()
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	s := models.NewShorten(short, req.URL)
	err = storage.GetStorage().AddShorten(*s)
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	res := models.Response{
		Result: fmt.Sprintf("%s/%s", config.Config.Base, short),
	}

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)

	res := models.BadResponse{
		Code:    statusCode,
		Message: message,
	}

	if err := encoder.Encode(res); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
	}
}

func GetRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(
		compresser.CompressMiddleware,
		logger.LogMiddleware,
	)

	router.Route("/api", func(r chi.Router) {
		r.Post("/shorten", apiCreateShortHandler)
	})

	router.Get("/{short}", tryRedirectHandler)
	router.Post("/", createShortHandler)

	return router
}

func main() {
	config.InitConfig()

	log.Info().Msgf("[*] Application address: %s", config.Config.Base)
	log.Info().Msgf("[*] Base address: %s", config.Config.Base)

	if err := http.ListenAndServe(config.Config.App, GetRouter()); err != nil {
		panic(err)
	}
}
