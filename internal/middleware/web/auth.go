package web

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/auth"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.c
	"strings"
)

func Authorization(c flamego.Context, r flamego.Render) {
	token := c.Request().Header.Get("Authorization")
	if token == "" || strings.Index(token, "Bearer") != 0 {
		response.UnAuthorization(r)
		return
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	entry, err := auth.ParseToken(token)
	if err != nil {
		response.UnAuthorization(r)
		return
	}
	c.Map(entry.Info)
}
