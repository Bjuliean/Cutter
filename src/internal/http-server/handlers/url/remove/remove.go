package remove

import (
	"net/http"
	resp "rapi/rapi/src/internal/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const ferr = "handlers.url.remove.New"

		log := log.With(
			slog.String("ferr", ferr),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("invalid request")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Info("failed to delete alias")
			render.JSON(w, r, resp.Error("failed to delete alias"))
			return
		}

		log.Info("alias deleted")
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
