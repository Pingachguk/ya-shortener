package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShortHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		data string
		want want
	}{
		{
			name: "#1 Redirected",
			data: "qwerty",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name: "#2 Not Found",
			data: "qweasdzxc",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	urls := UrlStorage{
		"qwerty": "https://praktikum.yandex.ru",
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s", test.data), nil)
			w := httptest.NewRecorder()
			handler := getShortHandler(urls)
			handler(w, r)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}

func TestCreateShortHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		data string
		want want
	}{
		{
			name: "#1 Created",
			data: "https://praktikum.yandex.ru",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "#2 Bad Request",
			data: "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	urls := make(UrlStorage)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(test.data)))
			w := httptest.NewRecorder()
			handler := createShortHandler(urls)
			handler(w, r)

			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}
