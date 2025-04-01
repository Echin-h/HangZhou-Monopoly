package dao

import (
	"errors"

	"github.com/Echin-h/HangZhou-Monopoly/config"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/database"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Init(db *gorm.DB) error {
	if database.GetDb(config.GetConfig().Databases[0].Key) == nil {
		return errors.New("database not found")
	}
	DB = db
	return DB.AutoMigrate(&model.User{})
}
