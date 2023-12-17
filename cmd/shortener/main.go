package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pingachguk/ya-shortener/config"
)

type URLStorage map[string]string

var urls URLStorage

func tryRedirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(123)
	id := chi.URLParam(r, "id")
	url := urls[id]

	if len(url) == 0 {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func createShortHandler(w http.ResponseWriter, r *http.Request) {
	if urls == nil {
		urls = make(URLStorage)
	}

	b := make([]byte, r.ContentLength)
	r.Body.Read(b)
	url := string(b)

	if len(url) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	short := fmt.Sprint(len(urls) + 1)
	urls[short] = url
	res := fmt.Sprintf("%s/%s", cfg.Base, short)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func GetRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/{id}", tryRedirectHandler)
	router.Post("/", createShortHandler)

	return router
}

var cfg config.Config

func main() {
	cfg = config.New()

	fmt.Printf("[*] Application address: %s\n", cfg.App)
	fmt.Printf("[*] Base address: %s\n", cfg.Base)

	if err := http.ListenAndServe(cfg.App, GetRouter()); err != nil {
		panic(err)
	}
}
