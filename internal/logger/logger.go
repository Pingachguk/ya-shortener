package logger

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Info().Msgf("[%s] %d %s %s", r.RemoteAddr, r.Response.StatusCode, r.Method, r.RequestURI)
	})
}
