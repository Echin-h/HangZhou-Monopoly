package dao

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/ping/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/stdao"
	"gorm.io/gorm"
)

var (
	Ping    *gorm.DB
	PingDto = stdao.Create(&model.Ping{})
)

func Init(db *gorm.DB) error {
	err := PingDto.Init(db)
	if err != nil {
		return err
	}
	Ping = db
	return AutoMigrate()
}

func AutoMigrate() error {
	return Ping.AutoMigrate(&model.Ping{})
}
