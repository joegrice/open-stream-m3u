package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joe/open-stream-m3u/internal/config"
	"github.com/joe/open-stream-m3u/internal/server"
)

//go:embed all:web
var webFS embed.FS

func main() {
	cfg := config.Load()

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		logger.Error("Failed to load web assets", "error", err)
		os.Exit(1)
	}

	srv := server.New(cfg, logger, webContent)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Open Stream M3U started", "port", cfg.Port)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down...")
}
