package delete

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/raisultan/url-shortener/internal/lib/api/response"
	"github.com/raisultan/url-shortener/internal/lib/logger/sl"
	"github.com/raisultan/url-shortener/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlDeleterStorage
type UrlDeleterStorage interface {
	DeleteUrl(ctx context.Context, alias string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlDeleterCache
type UrlDeleterCache interface {
	DeleteUrl(ctx context.Context, alias string) error
}

func New(
	log *slog.Logger,
	urlDeleterStorage UrlDeleterStorage,
	urlDeleterCache UrlDeleterCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		err := urlDeleterStorage.DeleteUrl(r.Context(), alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("alias not found", slog.String("alias", alias))
			render.JSON(w, r, response.Error("alias not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to delete url"))
			return
		}

		err = urlDeleterCache.DeleteUrl(r.Context(), alias)
		if err != nil {
			log.Error("failed to delete url from cache", sl.Err(err))
		}

		log.Info("url deleted")
		render.JSON(w, r, response.OK())
	}
}
