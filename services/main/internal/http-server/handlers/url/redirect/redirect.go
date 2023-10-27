package redirect

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/raisultan/url-shortener/lib/api/response"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/main/internal/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

const (
	urlNotFoundMessage   = "url not found"
	internalErrorMessage = "internal error"
)

type UrlGetterStorage interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

type UrlGetterCache interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

type AnalyticsTracker interface {
	TrackClickEvent(
		r *http.Request,
		alias string,
		latency time.Duration,
		errMessage string,
	) error
}

func New(
	log *slog.Logger,
	urlGetterStorage UrlGetterStorage,
	urlGetterCache UrlGetterCache,
	analyticsTracker AnalyticsTracker,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		startTime := time.Now()
		alias := chi.URLParam(r, "alias")
		resUrl, errMessage := getUrl(r.Context(), log, urlGetterCache, urlGetterStorage, alias)

		latency := time.Since(startTime)
		err := analyticsTracker.TrackClickEvent(r, alias, latency, errMessage)
		if err != nil {
			log.Error("failed to send click analytics event", sl.Err(err))
		} else {
			log.Info("click event sent to analytics storage")
		}

		if errMessage != "" {
			render.JSON(w, r, response.Error(errMessage))
		} else {
			http.Redirect(w, r, resUrl, http.StatusFound)
		}
	}
}

func getUrl(
	ctx context.Context,
	log *slog.Logger,
	urlGetterCache UrlGetterCache,
	urlGetterStorage UrlGetterStorage,
	alias string,
) (string, string) {
	resUrl, err := urlGetterCache.GetUrl(ctx, alias)
	if err == nil {
		log.Info("got url from cache", slog.String("url", resUrl))
		return resUrl, ""
	}

	log.Info("url not found in cache, checking storage", "alias", alias)
	resUrl, err = urlGetterStorage.GetUrl(ctx, alias)

	if errors.Is(err, storage.ErrUrlNotFound) {
		log.Info("url not found in storage", "alias", alias)
		return "", urlNotFoundMessage
	}

	if err != nil {
		log.Error("failed to get url from storage", sl.Err(err))
		return "", internalErrorMessage
	}

	log.Info("got url from storage", slog.String("url", resUrl))
	return resUrl, ""
}
