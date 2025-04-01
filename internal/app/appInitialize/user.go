package appInitialize

import "github.com/Echin-h/HangZhou-Monopoly/internal/app/user"

func init() {
	apps = append(apps, &user.User{Name: "user module"})
}
