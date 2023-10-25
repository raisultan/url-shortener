package redirect

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	"github.com/raisultan/url-shortener/lib/api/response"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/main/internal/analytics"
	"github.com/raisultan/url-shortener/services/main/internal/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=AnalyticsTracker
type AnalyticsTracker interface {
	TrackClickEvent(event analytics.ClickEvent) error
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
		resUrl, err := urlGetterCache.GetUrl(r.Context(), alias)
		if err != nil {
			log.Info("url not found in cache, checking storage", "alias", alias)
			resUrl, err = urlGetterStorage.GetUrl(r.Context(), alias)
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Info("url not found", "alias", alias)
				render.JSON(w, r, response.Error("url not found"))
				return
			}
			if err != nil {
				log.Error("failed to get url from storage", sl.Err(err))
				render.JSON(w, r, response.Error("internal error"))
				return
			} else {
				log.Info("got url from storage", slog.String("url", resUrl))
			}
		} else {
			log.Info("got url from cache", slog.String("url", resUrl))
		}

		latency := time.Since(startTime)
		event := analytics.ClickEvent{
			URLAlias:  alias,
			Timestamp: time.Now(),
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
			Referrer:  r.Referer(),
			Latency:   latency,
		}

		err = analyticsTracker.TrackClickEvent(event)
		if err != nil {
			log.Error("failed to send click analytics event", sl.Err(err))
		} else {
			log.Info("sent click analytics event", slog.Any("click-event", event))
		}

		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
