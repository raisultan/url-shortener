package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	middlewareLogger "github.com/raisultan/url-shortener/lib/http-server/middleware/logger"
	"github.com/raisultan/url-shortener/lib/logger"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/main/internal/alias"
	"github.com/raisultan/url-shortener/services/main/internal/analytics/clickhouse"
	"github.com/raisultan/url-shortener/services/main/internal/cache/redis"
	"github.com/raisultan/url-shortener/services/main/internal/config"
	"github.com/raisultan/url-shortener/services/main/internal/http-server/handlers/url/delete"
	"github.com/raisultan/url-shortener/services/main/internal/http-server/handlers/url/redirect"
	"github.com/raisultan/url-shortener/services/main/internal/http-server/handlers/url/save"
	"github.com/raisultan/url-shortener/services/main/internal/storage/mongo"
	"github.com/raisultan/url-shortener/services/main/internal/storage/sqlite"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Storage interface {
	Close(_ context.Context, log *slog.Logger)
	SaveUrl(
		_ context.Context,
		urlToSave string,
		alias string,
	) error
	GetUrl(_ context.Context, alias string) (string, error)
	DeleteUrl(_ context.Context, alias string) error
}

func main() {
	cfg := config.MustLoadConfig()
	log := logger.SetupLogger(cfg.Env)

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

	cache, err := redis.New(cfg.Cache, ctx)
	if err != nil {
		log.Error("failed to initialize cache", sl.Err(err))
		os.Exit(1)
	}
	defer cache.Close(log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middlewareLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	agc := alias.NewAliasGeneratorClient(cfg.AliasGenerator)
	analyticsTracker := clickhouse.NewClickHouseAnalyticsTracker()

	router.Post("/url", save.New(log, storage, cache, agc))
	router.Get("/{alias}", redirect.New(log, storage, cache, analyticsTracker))
	router.Delete("/{alias}", delete.New(log, storage, cache))

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
