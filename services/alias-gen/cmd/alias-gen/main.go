package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	middlewareLogger "github.com/raisultan/url-shortener/lib/http-server/middleware/logger"
	"github.com/raisultan/url-shortener/lib/logger"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/alias-gen/internal/config"
	"github.com/raisultan/url-shortener/services/alias-gen/internal/http-server/handlers/alias/generate"
	"github.com/raisultan/url-shortener/services/alias-gen/internal/storage/postgres"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting alias-generator", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	defer storage.Close(log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middlewareLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/alias", generate.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.HttpServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done

	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.HttpServer.CtxTimeout,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}
