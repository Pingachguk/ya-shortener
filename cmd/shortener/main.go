package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLStorage map[string]string

var urls URLStorage

func getShortHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url := urls[id]

	if len(url) == 0 {
		http.NotFound(w, r)
	} else {
		w.WriteHeader(http.StatusTemporaryRedirect)
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
	res := fmt.Sprintf("http://%s/%s", r.Host, short)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func GetRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/{id}", getShortHandler)
	router.Post("/", createShortHandler)

	return router
}

func main() {
	if err := http.ListenAndServe(":8080", GetRouter()); err != nil {
		panic(err)
	}
}
