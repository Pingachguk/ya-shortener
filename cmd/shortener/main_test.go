package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/pingachguk/ya-shortener/config"
	"github.com/rs/zerolog/log"
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

	srv := createTestServer()
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

	srv := createTestServer()
	defer srv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = fmt.Sprintf("%s%s", srv.URL, "/")
			req.SetBody(test.data)

			res, err := req.Send()

			require.NoError(t, err, "Error HTTP Request")
			assert.Equal(t, test.want.statusCode, res.StatusCode())
		})
	}
}

func TestApiCreateShortHandler(t *testing.T) {
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
			data: `{"url": "https://praktikum.yandex.ru"}`,
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

	srv := createTestServer()
	defer srv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = fmt.Sprintf("%s%s", srv.URL, "/api/shorten")
			req.SetHeader("Content-Type", "application/json")
			req.SetBody(test.data)

			res, err := req.Send()
			require.NoError(t, err, "Error HTTP Request")
			assert.Equal(t, test.want.statusCode, res.StatusCode())
		})
	}
}

func TestCompress(t *testing.T) {
	requestBody := `{"url": "https://praktikum.yandex.ru"}`

	srv := createTestServer()
	defer srv.Close()

	t.Run("Send data compressed", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = fmt.Sprintf("%s%s", srv.URL, "/api/shorten")
		req.SetHeader("Content-Encoding", "gzip")
		req.SetHeader("Content-Type", "application/json")

		buf := bytes.NewBuffer(nil)
		zw := gzip.NewWriter(buf)
		_, err := zw.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zw.Close()
		require.NoError(t, err)

		req.SetBody(buf)

		res, err := req.Send()

		require.NoError(t, err, "Error HTTP Request")
		assert.Equal(t, http.StatusCreated, res.StatusCode())
	})

	t.Run("Decompress data", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.URL = fmt.Sprintf("%s%s", srv.URL, "/api/shorten")
		req.SetHeader("Accept-Encoding", "gzip")
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(requestBody)
		req.SetDoNotParseResponse(true)

		res, err := req.Send()

		require.NoError(t, err, "Error HTTP Request")

		assert.Equal(t, http.StatusCreated, res.StatusCode(), "Неожиданный кож ответа")
		contentEncoding := res.Header().Get("Content-Encoding")
		hasGzip := strings.Contains(contentEncoding, "gzip")
		assert.True(t, hasGzip, "Заголовок Content-Encoding не содержит значение gzip")

		zr, err := gzip.NewReader(res.RawResponse.Body)
		require.NoError(t, err, "Ошибка инициализации gzip.Reader")

		result, err := io.ReadAll(zr)
		require.NoError(t, err, "Ошибка чтения gzip.Reader")
		log.Info().Msgf("%s", string(result))
	})
}

func createTestServer() *httptest.Server {
	config.InitConfig()
	srv := httptest.NewServer(GetRouter())

	return srv
}
