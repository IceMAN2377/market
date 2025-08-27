package main

import (
	"github.com/IceMAN2377/market/app"
	"github.com/IceMAN2377/market/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"strings"
)

func main() {
	cfg := config.NewConfig()

	// Настройка логгера в зависимости от уровня
	var logLevel slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	app := app.NewApp(cfg, logger)

	go func() {
		logger.Info("Starting subscription service",
			slog.Int("port", cfg.HttpPort),
			slog.String("log_level", cfg.LogLevel))

		if err := app.Run(); err != nil {
			logger.Error("Failed to start HTTP server", "error", err)
			panic("failed to start http server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logger.Info("Shutting down subscription service")
}
