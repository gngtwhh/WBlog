package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	// 假设你有一个简单的生成随机字符串的工具，如果没有，暂时用 time.Now().UnixNano() 代替
	// "github.com/google/uuid"
)

type contextKey string

const LoggerKey contextKey = "logger"

// RequestLogger record request log and inject Request-ID
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := time.Now().Format("20060102150405.000000")

			// child logger
			reqLogger := logger.With(
				slog.String("req_id", reqID),
			)
			ctx := context.WithValue(r.Context(), LoggerKey, reqLogger)

			next.ServeHTTP(w, r.WithContext(ctx))

			duration := time.Since(start)
			reqLogger.Info("http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("ip", r.RemoteAddr),
				slog.Int64("cost_ms", duration.Milliseconds()),
			)
		})
	}
}

// GetLogger gets Logger instance from context，if
// failed, returns default logger
func GetLogger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
