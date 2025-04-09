package router

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/handler"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/web"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func AppUserInit(e *flamego.Flame) {
	e.Group("/user/v1", func() {
		//e.Get("", handler.HandlePing)
		e.Group("/login", func() {
			//e.Post("/general", binding.JSON(dto.GeneralLoginRequest{}), handler.HandleGeneralLogin)
			e.Post("/test", handler.HandleTest)
			e.Post("/wechat", binding.JSON(dto.WechatLoginRequest{}), handler.HandleWechatLogin)
		})

		e.Group("/register", func() {
			e.Post("/wechat", binding.JSON(dto.WechatRegisterRequest{}), handler.HandleRegister)
		})
		e.Group("", func() {
			e.Get("/me", handler.HandleGetMyInfo)
		}, web.Authorization)
		e.Group("", func() {
			e.Post("/refresh", binding.JSON(dto.RefreshTokenRequest{}), handler.HandleRefreshToken) // 更新Token
		})
	})
}
