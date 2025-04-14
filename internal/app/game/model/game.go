package model

import (
	"math/rand"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GameState int

const (
	GameStatePrepared GameState = iota
	GameStateStarted
	GameStateFinished

	MaxTeamMember  = 4
	MaxPlayerCount = 16
)

type Dice struct {
	Value int `json:"value"`
}

type Game struct {
	ID          string    `gorm:"primaryKey;type:char(26);"`
	Code        int       `gorm:"code;index"`   // 游戏码
	Sponsor     string    `gorm:"sponsor"`      // 发起人
	State       GameState `gorm:"type:int"`     // 游戏状态
	PlayerCount int       `gorm:"player_count"` // 玩家数量 20 个
	TeamCount   int       `gorm:"team_count"`   // 队伍数量 4支
}

// 正确的GORM钩子签名
func (g *Game) BeforeCreate(tx *gorm.DB) error {
	// 生成ULID
	g.ID = ulid.Make().String()

	// 生成6位随机码
	g.Code = rand.Intn(900000) + 100000 // 100000-999999

	// 设置默认状态
	if g.State == 0 {
		g.State = GameStatePrepared
	}

	return nil
}
