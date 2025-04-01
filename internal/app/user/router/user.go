package router

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/handler"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppUserInit(e *flamego.Flame) {
	e.Get("/user/v1", handler.HandleExampleGet)
	e.Post("/user/v1", binding.JSON(dto.ExamplePost{}), handler.HandleExamplePost)
}
