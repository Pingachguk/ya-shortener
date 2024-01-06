package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingachguk/ya-shortener/config"
	"github.com/pingachguk/ya-shortener/internal/logger"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
)

var urls map[string]string

func tryRedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url, ok := urls[id]

	if ok {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func createShortHandler(w http.ResponseWriter, r *http.Request) {
	if urls == nil {
		urls = make(map[string]string)
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	} else if len(body) == 0 {
		http.Error(w, "Bad reuqest data: empty body", http.StatusBadRequest)
		return
	}

	url := string(body)
	short, err := shortid.GetDefault().Generate()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	urls[short] = url
	res := fmt.Sprintf("%s/%s", cfg.Base, short)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func GetRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(logger.LogMiddleware)

	router.Get("/{id}", tryRedirectHandler)
	router.Post("/", createShortHandler)

	return router
}

var cfg config.Config

func main() {
	cfg = config.New()

	log.Info().Msgf("[*] Application address: %s", cfg.App)
	log.Info().Msgf("[*] Base address: %s", cfg.Base)

	if err := http.ListenAndServe(cfg.App, GetRouter()); err != nil {
		panic(err)
	}
}
