package service

import (
	"math/rand"
	"sync"
	"time"

	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/model"
)

type GameService struct {
	Teams       []*model.Team
	Stations    map[string]*model.Station
	Lines       map[model.LineType]*model.MetroLine
	CurrentTurn int
	GameState   model.GameState
	StartTime   time.Time
	Dice        model.Dice
	mu          sync.Mutex
}

func NewGameService() *GameService {
	game := &GameService{
		Teams:     make([]*model.Team, 0),
		Stations:  make(map[string]*model.Station),
		Lines:     make(map[model.LineType]*model.MetroLine),
		GameState: model.GameStatePrepared,
	}

	game.initMetroLines()
	return game
}

// 初始化地铁线路
func (g *GameService) initMetroLines() {
	// 从配置文件加载站点数据
	stationData := loadStationConfig("config/stations.json")

	// 初始化线路
	g.Lines = map[model.LineType]*model.MetroLine{
		model.Line1:  {Type: model.Line1, Name: "1号线"},
		model.Line2:  {Type: model.Line2, Name: "2号线"},
		model.Line3:  {Type: model.Line3, Name: "3号线"},
		model.Line4:  {Type: model.Line4, Name: "4号线"},
		model.Line5:  {Type: model.Line5, Name: "5号线"},
		model.Line19: {Type: model.Line19, Name: "19号线"},
	}

	// 添加站点到线路
	for _, config := range stationData.Stations {
		station := &model.Station{
			ID:             config.ID,
			Name:           config.Name,
			Type:           model.StationType(config.Type),
			Price:          config.Price,
			CurrentStation: config.CurrentStation,
			Line:           model.LineType(config.Line),
		}
		g.Stations[config.ID] = station
		g.Lines[model.LineType(config.Line)].Stations = append(
			g.Lines[model.LineType(config.Line)].Stations, station)
	}

	// 连接同一线路的站点
	for _, line := range g.Lines {
		for i := 0; i < len(line.Stations)-1; i++ {
			line.Stations[i].Next = line.Stations[i+1]
			line.Stations[i+1].Prev = line.Stations[i]
		}
	}

	// 连接换乘站
	for _, s1 := range g.Stations {
		if s1.Type == model.StationTypeTransfer {
			for _, s2 := range g.Stations {
				if s2.Type == model.StationTypeTransfer && s1.ID != s2.ID &&
					s1.CurrentStation.X == s2.CurrentStation.X && s1.CurrentStation.Y == s2.CurrentStation.Y {
					connectStations(s1, s2)
				}
			}
		}
	}
}

func connectStations(s1, s2 *model.Station) {
	s1.ConnectedTo = append(s1.ConnectedTo, s2)
	s2.ConnectedTo = append(s2.ConnectedTo, s1)
}

// 游戏核心逻辑方法
func (g *GameService) AddTeam(team *model.Team) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 设置初始位置
	if len(g.Stations) > 0 {
		for _, station := range g.Stations {
			team.CurrentStation = station
			break
		}
	}

	g.Teams = append(g.Teams, team)
}

func (g *GameService) RollDice(teamID string) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.GameState != model.GameStateStarted {
		return 0
	}

	// 掷骰子
	value := rand.Intn(6) + 1
	g.Dice.Value = value

	// 移动当前队伍
	team := g.findTeam(teamID)
	g.moveTeamBySteps(team, value)

	// 切换到下一个队伍
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Teams)

	return value
}

func (g *GameService) moveTeamBySteps(team *model.Team, steps int) {
	// 计算新的位置
	for i := 0; i < steps; i++ {
		if team.CurrentStation.Next != nil {
			team.CurrentStation = team.CurrentStation.Next
		} else {
			team.CurrentStation = g.Lines[team.CurrentStation.Line].Stations[0]
		}
	}

	// 检查是否需要购买站点
	if team.CurrentStation.Owner == nil {
		g.buyStation(team)
	}
}

func (g *GameService) buyStation(team *model.Team) {
	// 购买站点逻辑
	if team.Balance >= team.CurrentStation.Price {
		team.Balance -= team.CurrentStation.Price
		team.Stations = append(team.Stations, team.CurrentStation)
		team.CurrentStation.Owner = team
	}
}

func (g *GameService) findTeam(teamID string) *model.Team {
	for _, team := range g.Teams {
		if team.ID == teamID {
			return team
		}
	}
	return nil
}
