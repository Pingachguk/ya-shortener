package main

import (
	"fmt"
	"net/http"
	"strings"
)

var urls map[string]string

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, id, _ := strings.Cut(r.RequestURI, "/")
		url := getShort(id)
		if len(url) == 0 {
			http.NotFound(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
	case http.MethodPost:
		b := make([]byte, r.ContentLength)
		r.Body.Read(b)
		short := createShort(string(b))
		res := fmt.Sprintf("http://%s/%s", r.Host, short)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(res))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getShort(id string) string {
	return urls[id]
}

func createShort(url string) string {
	res := fmt.Sprint(len(urls) + 1)
	urls[res] = url

	return res
}

func main() {
	urls = make(map[string]string)
	mux := http.NewServeMux()

	mux.HandleFunc("/", mainHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
