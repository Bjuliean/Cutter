package main

import (
	"fmt"
	"os"
	"rapi/rapi/src/internal/config"
	_"rapi/rapi/src/internal/storage"
	"rapi/rapi/src/internal/storage/sqlite"

	"golang.org/x/exp/slog"
)

const(
	envLocal="local"
	envDev="dev"
	envProd="prod"
)

func main() {
	//init cfg (cleanenv)
	//init logger (slog)
	//init storage (sqlite)
	//init router (chi, chi render)
	//run server

	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := createLogger(cfg.Env)
	log.Info("ABOBA INFO")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}
	storage.AAA()
}

func createLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug,}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug,}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo,}))
	}
	return logger
}