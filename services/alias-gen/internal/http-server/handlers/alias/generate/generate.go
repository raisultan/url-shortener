package generate

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/raisultan/url-shortener/lib/api/response"
	"github.com/raisultan/url-shortener/lib/logger/sl"
	"github.com/raisultan/url-shortener/services/alias-gen/internal/generator"
	"golang.org/x/exp/slog"
	"net/http"
)

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type CounterIncrementer interface {
	IncrementCounter() (int64, error)
}

func New(log *slog.Logger, counterIncrementer CounterIncrementer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.alias.generate.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		count, err := counterIncrementer.IncrementCounter()
		if err != nil {
			log.Error("failed to increment counter", sl.Err(err))
			render.JSON(w, r, response.Error("failed to increment counter"))
			return
		}

		alias := generator.GenerateAlias(count)

		log.Info("alias generated", slog.String("alias", alias))
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
