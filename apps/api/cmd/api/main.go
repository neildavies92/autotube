package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/neildavies92/autotube/apps/api/internal/config"
	"github.com/neildavies92/autotube/apps/api/internal/httpserver"
	"github.com/neildavies92/autotube/apps/api/internal/providers"
	"github.com/neildavies92/autotube/apps/api/internal/youtube"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load(config.Options{})
	if err != nil {
		logger.Error("configuration failed", "error", err)
		os.Exit(1)
	}

	var youtubeProvider providers.YouTubeProvider
	if cfg.Providers.YouTubeAPIKey != "" {
		youtubeProvider = youtube.NewClient(cfg.Providers.YouTubeAPIKey, nil)
		logger.Info("youtube provider configured")
	} else {
		logger.Warn("YOUTUBE_API_KEY not set — /research/youtube endpoint disabled")
	}

	server := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           httpserver.New(logger, youtubeProvider),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errs := make(chan error, 1)
	go func() {
		logger.Info("api server starting", "addr", server.Addr)
		errs <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errs:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", "error", err)
			os.Exit(1)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("api server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("api server stopped")
}
