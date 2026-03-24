package main

import (
	"URL-shortener/internal/config"
	"URL-shortener/internal/lib/logger/sl"
	"URL-shortener/internal/storage/postgres"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Starting URL-shortener", slog.String("env", cfg.Env))
	log.Debug("Dbg mess")

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("Unable to connect to database", sl.Err(err))
		os.Exit(1)
	}
	_ = storage
	fmt.Println(cfg)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
