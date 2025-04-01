package handler

import (
	"fmt"

	"github.com/Echin-h/HangZhou-Monopoly/config"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/flamego"
)

func HandleExampleGet(c flamego.Context, r flamego.Render) {
	data := struct {
		UA         string
		Host       string
		Method     string
		Proto      string
		RemoteAddr string
		Message    string
	}{
		UA:         c.Request().UserAgent(),
		Host:       c.Request().Host,
		Method:     c.Request().Method,
		Proto:      c.Request().Proto,
		RemoteAddr: c.Request().RemoteAddr,
		Message:    fmt.Sprintf("Welcome to %s, version %s.", config.GetConfig().ProgramName, config.GetConfig().VERSION),
	}
	response.HTTPSuccess(r, data)
}

func HandleExamplePost(r flamego.Render, req dto.ExamplePost) {
	response.HTTPSuccess(r, req.Msg)
}
