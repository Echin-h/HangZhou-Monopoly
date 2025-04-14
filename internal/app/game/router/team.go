package router

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/handler"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/web"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppTeamInit(e *flamego.Flame) {
	e.Group("/team/v1", func() {
		e.Post("/team", binding.JSON(dto.TeamCreateRequest{}), handler.HandleTeamJoin)
	}, web.Authorization)
}
