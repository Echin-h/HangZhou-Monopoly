package model

import (
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Team struct {
	ID        string `gorm:"primaryKey;type:char(16);"` // 游戏码-1|2|3|4
	GameID    string `gorm:"game_id;type:char(16);index"`
	Name      string `gorm:"name"`
	Balance   int    `gorm:"balance;comment:金额"`
	Leader    string `gorm:"type:char(28);"`
	Count     int    `gorm:"count;comment:队伍人数"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt

	Line           LineType   `json:"line" gorm:"-"`           // 线路
	CurrentStation *Station   `json:"currentStation" gorm:"-"` // 当前站点
	Stations       []*Station `json:"stations" gorm:"-"`       // 拥有的站点数量
}

type TeamUser struct {
	TeamID    string `gorm:"team_id;type:char(28);index"`
	UserID    string `gorm:"user_id;type:char(28);index"`
	GameID    string `gorm:"game_id;type:char(16);index"`
	DeletedAt gorm.DeletedAt
}

func (t *Team) BeforeCreate() string {
	t.ID = ulid.Make().String()
	t.Name = "team-" + randomString(5)
	return t.ID
}

// randomString 函数用于生成指定长度的随机字符串
func randomString(length int) string {
	var sb strings.Builder
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		// 从字符集中随机选取一个字符
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}
	return sb.String()
}
