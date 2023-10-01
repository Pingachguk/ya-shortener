package main

import (
	"fmt"
	"net/http"
	"strings"
)

var urls map[string]string

func router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getShortHandler(w, r)
	case http.MethodPost:
		createShortHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getShortHandler(w http.ResponseWriter, r *http.Request) {
	_, id, _ := strings.Cut(r.RequestURI, "/")
	url := urls[id]
	if len(url) == 0 {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func createShortHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, r.ContentLength)
	r.Body.Read(b)
	url := string(b)

	short := fmt.Sprint(len(urls) + 1)
	urls[short] = url
	res := fmt.Sprintf("http://%s/%s", r.Host, short)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func main() {
	urls = make(map[string]string)
	mux := http.NewServeMux()

	mux.HandleFunc("/", router)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
