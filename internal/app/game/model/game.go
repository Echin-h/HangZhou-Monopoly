package model

import (
	"math/rand"

	"github.com/oklog/ulid/v2"
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
	ID          string    `gorm:"primaryKey;type:char(16);"`
	Code        int       `gorm:"code;index"`   // 游戏码
	Sponsor     string    `gorm:"sponsor"`      // 发起人
	State       GameState `gorm:"type:int"`     // 游戏状态
	PlayerCount int       `gorm:"player_count"` // 玩家数量 20 个
	TeamCount   int       `gorm:"team_count"`   // 队伍数量 4支
}

func (g *Game) BeforeCreate() string {
	g.ID = ulid.Make().String()
	g.Code = rand.Intn(900000) + 100000
	return g.ID
}
