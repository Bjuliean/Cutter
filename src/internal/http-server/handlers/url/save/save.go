package save

import (
	//"errors"
	"fmt"
	"net/http"
	resp "rapi/rapi/src/internal/api/response"
	"rapi/rapi/src/internal/random"
	//"rapi/rapi/src/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"golang.org/x/exp/slog"
	//"golang.org/x/text/number"
)

const aliasLength = 10

type URLSaver interface {
	SaveURL(newUrl string, alias string) error
}

type Request struct {
	URL		string		`json:"url" validate:"required,url"`
	Alias	string		`json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias	string		`json:"alias,omitempty"`
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const ferr = "handlers.url.save.New"

		log := log.With(
			slog.String("ferr", ferr),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failder to decode request body")
			render.JSON(w, r, resp.Error("failed to decode request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request")
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.RandomString(aliasLength)
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error(fmt.Sprintf("error: %s: %s", ferr, err.Error()))
			render.JSON(w, r, resp.Error(fmt.Sprintf("error: %s: %s", ferr, err.Error())))
			return
		}
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias: alias,
		})
	}
}