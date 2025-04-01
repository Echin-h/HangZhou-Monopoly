package app

import (
	"errors"

	"github.com/Echin-h/HangZhou-Monopoly/config"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/database"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/kernel"
)

type Module interface {
	Info() string
	PreInit(*kernel.Engine) error
	Init(*kernel.Engine) error
	PostInit(*kernel.Engine) error
	Load(*kernel.Engine) error
	Start(*kernel.Engine) error

	OnConfigChange() func(*kernel.Engine) error

	mustEmbedUnimplementedModule()
}

type UnimplementedModule struct{}

func (*UnimplementedModule) Info() string {
	return "unimplementedModule"
}

func (*UnimplementedModule) PreInit(*kernel.Engine) error {
	db := database.GetDb(config.GetConfig().Databases[0].Key)
	if db == nil {
		return errors.New("module user's database is null")
	}
	return dao.Init(db)
}

func (*UnimplementedModule) Init(*kernel.Engine) error {
	return nil
}

func (*UnimplementedModule) PostInit(*kernel.Engine) error {
	return nil
}

func (*UnimplementedModule) Load(*kernel.Engine) error {
	return nil
}

func (*UnimplementedModule) Start(*kernel.Engine) error {
	return nil
}

func (*UnimplementedModule) OnConfigChange() func(*kernel.Engine) error {
	return func(engine *kernel.Engine) error {
		return nil
	}
}

func (*UnimplementedModule) mustEmbedUnimplementedModule() {}
