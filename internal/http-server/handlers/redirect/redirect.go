package redirect

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlgetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal/http-server/handlers/redirect/NEW"

		log = slog.With(slog.String("op", op))

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		resUrl, err := urlgetter.GetUrl(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		if err != nil {
			log.Error("failed to get url")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		log.Info("got url", slog.String("url", resUrl))

		// redirect to found url:

		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
