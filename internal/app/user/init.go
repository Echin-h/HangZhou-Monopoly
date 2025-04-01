package user

import (
	"context"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/router"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/kernel"

	"sync"
)

type User struct {
	Name string

	app.UnimplementedModule
}

func (p *User) Info() string {
	return p.Name
}

func (p *User) PreInit(engine *kernel.Engine) error {
	return nil
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
