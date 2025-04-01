package web

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/websocket"
	"github.com/flamego/flamego"
)

func InjectWebsocket(key string) flamego.Handler {
	return func(c flamego.Context) {
		c.Map(websocket.GetSocketManager(key))
	}
}
