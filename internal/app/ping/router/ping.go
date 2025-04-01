package router

import (
	"errors"

	"github.com/Echin-h/HangZhou-Monopoly/internal/app/ping/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/ping/handler"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppPingInit(e *flamego.Flame) {
	e.Get("/ping/v1", handler.HandleExampleGet)

	e.Get("/ping/v1/err", func(r flamego.Render) {
		response.HTTPFail(r, 500000, "test error", errors.New("this is err"))
	})

	e.Post("/ping/v1", binding.JSON(dto.ExamplePost{}), handler.HandlePing)
}
