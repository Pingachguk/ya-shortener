package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingachguk/ya-shortener/config"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/pingachguk/ya-shortener/internal/storage"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
)

func TryRedirectHandler(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	shorten, err := storage.GetStorage().GetByShort(context.Background(), short)
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	if shorten != nil {
		http.Redirect(w, r, shorten.OriginalURL, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func CreateShortHandler(w http.ResponseWriter, r *http.Request) {
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
	err = storage.GetStorage().AddShorten(context.Background(), *s)
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", config.Config.Base, short)))
}

func APICreateShortHandler(w http.ResponseWriter, r *http.Request) {
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
	err = storage.GetStorage().AddShorten(context.Background(), *s)
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

func PingDatabase(w http.ResponseWriter, r *http.Request) {
	database := storage.GetDatabaseStorage()
	err := database.Conn.Ping(context.Background())
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("")
		return
	}

	w.WriteHeader(http.StatusOK)
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
