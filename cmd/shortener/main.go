package main

import (
	"fmt"
	"net/http"
	"strings"
)

type UrlStorage map[string]string

var urls UrlStorage

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getShortHandler(urls)
	case http.MethodPost:
		createShortHandler(urls)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getShortHandler(urls UrlStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, id, _ := strings.Cut(r.RequestURI, "/")
		url := urls[id]
		if len(url) == 0 {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
	}
}

func createShortHandler(urls UrlStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	urls = make(UrlStorage)
	mux := http.NewServeMux()

	mux.HandleFunc("/", router)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
