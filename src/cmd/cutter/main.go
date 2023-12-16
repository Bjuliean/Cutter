package main

import (
	"fmt"
	"net/http"
	"os"
	"rapi/rapi/src/internal/config"
	"rapi/rapi/src/internal/http-server/handlers/url/save"
	"rapi/rapi/src/internal/http-server/middleware/logger"
	_ "rapi/rapi/src/internal/storage"
	"rapi/rapi/src/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	//log.Info("ABOBA INFO")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	// err = storage.SaveURL("https://google.com", "google")
	// if err != nil {
	// 	log.Error(err.Error())
	// }

	// res, err := storage.GetURL("google")
	// if err != nil {
	// 	log.Error("failed to init storage", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("res: %s\n", res)

	err = storage.DeleteURL("google")
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer) //panic defense
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}

	log.Error("server stopped")
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