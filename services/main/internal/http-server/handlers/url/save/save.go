package save

import (
	"context"
	"errors"
	"github.com/raisultan/url-shortener/lib/api/response"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/main/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlSaverStorage
type UrlSaverStorage interface {
	SaveUrl(
		ctx context.Context,
		urlToSave string,
		alias string,
	) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UrlSaverCache
type UrlSaverCache interface {
	SaveUrl(
		ctx context.Context,
		urlToSave string,
		alias string,
	) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=AliasGenerator
type AliasGenerator interface {
	GenerateAlias() (string, error)
}

func New(
	log *slog.Logger,
	urlSaverStorage UrlSaverStorage,
	urlSaverCache UrlSaverCache,
	aliasGenerator AliasGenerator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias, err = aliasGenerator.GenerateAlias()
			if err != nil {
				log.Error("failed to get alias", sl.Err(err))
				render.JSON(w, r, response.Error("failed to get alias"))
				return
			}
		}

		err = urlSaverStorage.SaveUrl(r.Context(), req.Url, alias)
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.Url))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		err = urlSaverCache.SaveUrl(r.Context(), req.Url, alias)
		if err != nil {
			log.Error("failed to add url to cache", sl.Err(err))
		}

		log.Info("url added")
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
