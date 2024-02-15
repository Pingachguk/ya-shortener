package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pingachguk/ya-shortener/internal/auth"
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
	shorten, err := storage.GetStorage().GetByShort(r.Context(), short)
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("error get by short for redirect")
		return
	}

	if shorten != nil {
		http.Redirect(w, r, shorten.OriginalURL, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func CreateShortHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(*r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("didn't get current user")
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("unread body")
		return
	} else if len(body) == 0 {
		http.Error(w, "Bad request data: empty body", http.StatusBadRequest)
		return
	}

	url := string(body)
	short, err := shortid.GetDefault().Generate()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad generating short id")
		return
	}

	fmt.Println(user)
	s := models.NewShorten(short, url, user.UUID)
	store := storage.GetStorage()
	err = store.AddShorten(r.Context(), *s)
	if errors.Is(err, storage.ErrUnique) {
		shorten, err := store.GetByURL(r.Context(), url)
		if err != nil {
			errorResponse(w, "Internal error", http.StatusInternalServerError)
			log.Error().Err(err).Msgf("bad get by url")
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf("%s/%s", config.Config.Base, shorten.ShortURL)))
		return
	} else if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("query error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", config.Config.Base, short)))
}

func APICreateShortHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(*r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("didn't get current user")
		return
	}

	var req models.ShortenRequest

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		errorResponse(w, "Bad Content-Type: need application/json", http.StatusNotAcceptable)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err == io.EOF {
		errorResponse(w, "Bad request data: empty body", http.StatusBadRequest)
		return
	} else if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad decoded request body")
		return
	}

	short, err := shortid.GetDefault().Generate()
	if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad generating short id")
		return
	}

	s := models.NewShorten(short, req.URL, user.UUID)
	store := storage.GetStorage()
	err = store.AddShorten(r.Context(), *s)
	encoder := json.NewEncoder(w)
	if errors.Is(err, storage.ErrUnique) {
		shorten, err := store.GetByURL(r.Context(), req.URL)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusConflict)
		res := models.ShortenResponse{
			Result: fmt.Sprintf("%s/%s", config.Config.Base, shorten.ShortURL),
		}
		if err := encoder.Encode(res); err != nil {
			errorResponse(w, "Internal error", http.StatusInternalServerError)
			log.Error().Err(err).Msgf("error encode response")
			return
		}
		return
	} else if err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("query error")
		return
	}

	res := models.ShortenResponse{
		Result: fmt.Sprintf("%s/%s", config.Config.Base, short),
	}

	w.WriteHeader(http.StatusCreated)
	if err := encoder.Encode(res); err != nil {
		errorResponse(w, "Internal error", http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad encode response")
		return
	}
}

func APIBatchCreateShortHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(*r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("didn't get current user")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	requestData := make([]models.BatchShortenRequest, 0)
	requestDecoder := json.NewDecoder(r.Body)

	if err := requestDecoder.Decode(&requestData); err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad decode request")
		return
	}

	responseData := make([]models.BatchShortenResponse, 0, len(requestData))
	shortens := make([]models.Shorten, 0, len(requestData))
	for _, v := range requestData {
		short, err := shortid.GetDefault().Generate()
		if err != nil {
			errorResponse(w, "Internal error", http.StatusInternalServerError)
			log.Error().Err(err).Msgf("bad generating short id")
			return
		}

		shorten := models.NewShorten(short, v.OriginalURL, user.UUID)
		shortens = append(shortens, *shorten)
		responseRow := models.BatchShortenResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", config.Config.Base, short),
		}
		responseData = append(responseData, responseRow)
	}

	err = storage.GetStorage().AddBatchShorten(r.Context(), shortens)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("batch error")
		return
	}

	w.WriteHeader(http.StatusCreated)

	responseEncoder := json.NewEncoder(w)
	if err := responseEncoder.Encode(responseData); err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad encode reponse")
		return
	}
}

func PingDatabase(w http.ResponseWriter, r *http.Request) {
	database := storage.GetDatabaseStorage()
	err := database.Conn.Ping(r.Context())
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("bad ping database")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUserURLS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := auth.GetCurrentUser(*r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("didn't get current user")
		return
	}

	if *user == (models.User{}) {
		err = auth.Authenticate(w)
		w.WriteHeader(http.StatusUnauthorized)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			log.Error().Err(err).Msgf("didn't get authenticate")
		}

		return
	}

	fmt.Println(user.UUID)

	store := storage.GetStorage()
	shortens, err := store.GetUserURLS(r.Context(), user.UUID)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		log.Error().Err(err).Msgf("didn't get user urls")
		return
	}

	if len(shortens) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		res := make([]models.UserURL, 0, len(shortens))
		for _, shorten := range shortens {
			userURL := &models.UserURL{
				OriginalURL: shorten.OriginalURL,
				ShortURL:    shorten.ShortURL,
			}
			res = append(res, *userURL)
		}

		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(res); err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			log.Error().Err(err).Msgf("bad create response for user urls")
			return
		}
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
		log.Error().Err(err).Msgf("bad encode response")
	}
}
