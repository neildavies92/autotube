package httpserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type rootResponse struct {
	Service string `json:"service"`
	Message string `json:"message"`
}

func New(logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleRoot)

	return recoverer(logger, requestLogger(logger, mux))
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, rootResponse{
		Service: "autotube-api",
		Message: "AutoTube API scaffold is running.",
	})
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)

		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(started).Milliseconds(),
		)
	})
}

func recoverer(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered",
					"method", r.Method,
					"path", r.URL.Path,
					"panic", recovered,
					"stack", string(debug.Stack()),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
