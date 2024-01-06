package logger

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	statusCode    int
	executionTime time.Duration
	http.ResponseWriter
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := loggingResponseWriter{
			statusCode:     0,
			executionTime:  0,
			ResponseWriter: w,
		}

		start := time.Now()
		next.ServeHTTP(&response, r)
		duration := time.Since(start)

		response.executionTime = duration

		log.Info().Msgf("[%s] %s %d %s (execution time: %s)", r.RemoteAddr, r.Method, response.statusCode, r.RequestURI, response.executionTime)
	})
}
