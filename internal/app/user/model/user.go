package model

import "time"

type User struct {
	OpenId    string `gorm:"type:char(28);primaryKey;not null"`
	Username  string `gorm:"type:char(255);not null"`
	CreatedAt time.Time
}
