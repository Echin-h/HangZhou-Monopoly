package game

import (
	"context"
	"sync"

	"github.com/Echin-h/HangZhou-Monopoly/internal/app"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/router"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/kernel"
)

type Game struct {
	Name string

	app.UnimplementedModule
}

func (p *Game) Info() string {
	return p.Name
}

func (p *Game) PreInit(engine *kernel.Engine) error {
	return nil
}

func (p *Game) Init(*kernel.Engine) error {
	return nil
}

func (p *Game) Load(engine *kernel.Engine) error {
	router.AppGameInit(engine.Fg)
	router.AppTeamInit(engine.Fg)
	return nil
}

func (p *Game) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (p *Game) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {
		return nil
	}
}
