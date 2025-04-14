package appInitialize

import "github.com/Echin-h/HangZhou-Monopoly/internal/app/game"

func init() {
	apps = append(apps, &game.Game{Name: "game module"})
}
