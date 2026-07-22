package api

import (
	"log/slog"
	"net/http"
	"time"
)

func LogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrappedWriter, r)

			logger.Info("Kyaa~ New HTTP Request!",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("user-agent", r.UserAgent()),
				slog.String("IP", r.RemoteAddr),
				slog.Int("status", wrappedWriter.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
