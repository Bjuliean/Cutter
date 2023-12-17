package redirect

import (
	"net/http"
	resp "rapi/rapi/src/internal/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const ferr = "handlers.redirect.New"

		log := log.With(
			slog.String("ferr", ferr),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			log.Info("url not found")
			render.JSON(w, r, resp.Error("url not found"))
			return
		}

		log.Info("url received", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
