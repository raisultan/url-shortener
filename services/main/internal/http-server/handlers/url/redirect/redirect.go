package redirect

import (
	"context"
	"errors"
	"github.com/raisultan/url-shortener/services/main/internal/lib/api/response"
	"github.com/raisultan/url-shortener/services/main/internal/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlGetterStorage
type UrlGetterStorage interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlGetterCache
type UrlGetterCache interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

func New(
	log *slog.Logger,
	urlGetterStorage UrlGetterStorage,
	urlGetterCache UrlGetterCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		resUrl, err := urlGetterCache.GetUrl(r.Context(), alias)
		if err != nil {
			log.Info("url not found in cache", "alias", alias)
		} else {
			log.Info("got url from cache", slog.String("url", resUrl))
			http.Redirect(w, r, resUrl, http.StatusFound)
			return
		}

		resUrl, err = urlGetterStorage.GetUrl(r.Context(), alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("got url from storage", slog.String("url", resUrl))
		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
