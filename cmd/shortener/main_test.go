package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	urls = map[string]string{
		"qwerty": "https://praktikum.yandex.ru",
	}

	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	client := srv.Client()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
			res, err := client.Get(srv.URL + fmt.Sprintf("/%s", test.data))
			require.NoError(t, err)
			defer res.Body.Close()

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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(test.data)))
			w := httptest.NewRecorder()
			createShortHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}
