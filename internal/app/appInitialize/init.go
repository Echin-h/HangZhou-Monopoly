package appInitialize

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app"
)

var (
	apps = make([]app.Module, 0)
)

func GetApps() []app.Module {
	return apps
}
