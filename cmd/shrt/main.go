package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http-server/handlers/redirect"
	"awesomeProject/internal/http-server/handlers/url/save"
	"awesomeProject/internal/lib/logger/ms"
	"awesomeProject/internal/storage/msql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	//"log/slog"
	"golang.org/x/exp/slog"
	"os"
)

// Константы для логгера ( разные случаи: локальный, прод, и разрабский)
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Mustload()

	log := setupLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := mysql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", ms.Err(err))
		os.Exit(1)
	}

	// Добавление в таблицу алиаса:

	//id, err := storage.SaveURL("url6", "al6")
	//if err != nil {
	//	log.Error("Faild to save url", ms.Err(err))
	//}
	//_ = id
	//log.Info("saved url", slog.Int64("id", id))

	_ = storage

	router := chi.NewRouter()
	//
	//router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	//router.Use(mw.New())

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("awesomeProject", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		router.Post("/url", save.New(log, storage))

	})

	router.Get("/{alias}", redirect.New(log, storage))

	//ЗАПУСКАЕМ СЕРВЕР !!!!

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Info("server started")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
