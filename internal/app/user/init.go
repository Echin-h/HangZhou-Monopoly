package user

import (
	"context"
	"sync"

	"github.com/Echin-h/HangZhou-Monopoly/internal/app"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/router"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/database"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/kernel"
)

type User struct {
	Name string

	app.UnimplementedModule
}

func (p *User) Info() string {
	return p.Name
}

func (p *User) PreInit(engine *kernel.Engine) error {
	db := database.GetDb("mysql")
	if db == nil {
		return nil
	}
	return dao.Init(db)
}

func (p *User) Init(*kernel.Engine) error {
	return nil
}

func (p *User) Load(engine *kernel.Engine) error {
	router.AppUserInit(engine.Fg)
	return nil
}

func (p *User) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *User) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {
		return nil
	}
}
