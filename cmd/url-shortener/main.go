package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raisultan/url-shortener/internal/config"
	"github.com/raisultan/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/raisultan/url-shortener/internal/http-server/handlers/url/save"
	"github.com/raisultan/url-shortener/internal/http-server/middleware/logger"
	"github.com/raisultan/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/raisultan/url-shortener/internal/lib/logger/sl"
	"github.com/raisultan/url-shortener/internal/storage/mongo"
	"github.com/raisultan/url-shortener/internal/storage/sqlite"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "production"
)

type Storage interface {
	Close(_ context.Context, log *slog.Logger)
	SaveUrl(
		_ context.Context,
		urlToSave string,
		alias string,
	) error
	GetUrl(_ context.Context, alias string) (string, error)
}

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.HttpServer.CtxTimeout,
	)
	defer cancel()

	storage, err := NewStorage(cfg, ctx)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}
	defer storage.Close(ctx, log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func NewStorage(cfg *config.Config, ctx context.Context) (Storage, error) {
	switch cfg.ActiveStorage {
	case "sqlite":
		return sqlite.New(
			cfg.Storages,
			ctx,
		)
	case "mongo":
		return mongo.New(
			cfg.Storages,
			ctx,
		)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.ActiveStorage)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
