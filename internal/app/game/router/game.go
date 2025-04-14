package router

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/handler"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/web"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppGameInit(e *flamego.Flame) {
	e.Group("/game/v1", func() {
		e.Post("/game", handler.HandleGameCreate)
		e.Post("/game/join", binding.JSON(dto.GameJoinRequest{}), handler.HandleGameJoin)
	}, web.Authorization)
}
