package save

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/ms"
	"awesomeProject/internal/lib/random"
	"awesomeProject/internal/storage"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

// обрабатываем джисон запрос в данную структуру
type Requests struct {
	URL   string `json:"url" validate:"required, url"`
	Alias string `json:"alias,omitempty"`
}

const aliaslength = 5

// кидаем джисон ответ из данной структуры
// omitempty  означает что если строчка пустая ее дропнет из джисона

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

// Интерфейс для сохранения запросов (берет базу (структура в папке msql))

type UrlSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op))

		var req Requests

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			// обработка ошибки с пустым запросом
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("Failed to decode json", ms.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		// Прописываем валидатор (Сам валидатор находится в internal/lib/api/response/response.go)

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", ms.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// ГЕНЕРАТОР АЛИАСОВ ЕСЛИ НЕТ КОНКРЕТНОГО
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliaslength)
		}
		id, err := urlSaver.SaveURL(
			req.URL,
			alias,
		)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", ms.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}
func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
