package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLStorage map[string]string

var urls URLStorage

func getShortHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "url")
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

	fmt.Println(r.Body)
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

func main() {
	router := chi.NewRouter()
	router.Get("/{url}", getShortHandler)
	router.Post("/", createShortHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
