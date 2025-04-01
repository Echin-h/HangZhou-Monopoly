package appInitialize

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/ping"
)

func init() {
	apps = append(apps, &ping.Ping{Name: "ping module"})
}
