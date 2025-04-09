package handler

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/auth"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/flamego"
)

func HandleTeamJoin(c flamego.Context, r flamego.Render, auth auth.Info) {
	var uid string
	if uid = auth.Uid; uid == "" {
		response.ErrorResponse(r, model.UnauthorizedErrorCode, "Unauthorized", "未授权")
		return
	}

	gameId := c.Param("game_id")
	if gameId == "" {
		response.ErrorResponse(r, model.ParamErrorCode, "Game ID empty", "游戏ID不能为空")
		return
	}
	teamId := c.Param("team_id")

	// 一个人只能加入一个队伍
	var tu []model.TeamUser
	if err := dao.DB.WithContext(c.Request().Context()).Where("user_id = ?", uid).Find(&tu).Error; err != nil {
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}
	if len(tu) > 0 {
		response.ErrorResponse(r, model.TeamAlreadyJoinErrorCode, "Team already joined", "已经加入队伍")
		return
	}

	// 游戏人数已满 | 游戏不存在
	var g model.Game
	if err := dao.DB.WithContext(c.Request().Context()).Where("id = ?", c.Param("game_id")).First(&g).Error; err != nil {
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}
	if g.PlayerCount == model.MaxPlayerCount {
		response.ErrorResponse(r, model.GameFullErrorCode, "Game full", "游戏人数已满")
		return
	}

	// 如果teamId为空表示创建队伍，不为空表示加入队伍
	if teamId == "" {
		// 有无队伍
		var ts []model.Team
		if err := dao.DB.WithContext(c.Request().Context()).Where("game_id = ?", gameId).Find(&ts).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
			return
		}

		if len(ts) > 0 {
			// 加入已有队伍
			joinTeamId := ""
			for _, t := range ts {
				if t.Count < model.MaxTeamMember {
					joinTeamId = t.ID
					break
				}
			}

			if joinTeamId == "" {
				response.ErrorResponse(r, model.TeamFullErrorCode, "Team full", "队伍已满")
				return
			}

			// 加入队伍
			if err := dao.DB.WithContext(c.Request().Context()).Create(&model.TeamUser{
				TeamID: joinTeamId,
				UserID: uid,
				GameID: gameId,
			}).Error; err != nil {
				response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
				return
			}

			if err := dao.DB.WithContext(c.Request().Context()).Model(&model.Team{}).Where("id = ?", joinTeamId).UpdateColumns(map[string]interface{}{
				"count": g.PlayerCount + 1,
			}).Error; err != nil {
				response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
				return
			}

			if err := dao.DB.WithContext(c.Request().Context()).Model(&model.Game{}).Where("id = ?", gameId).UpdateColumns(map[string]interface{}{
				"player_count": g.PlayerCount + 1,
				"team_count":   g.TeamCount,
			}).Error; err != nil {
				response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
				return
			}
		} else if len(ts) == 0 && g.PlayerCount < model.MaxPlayerCount { // 如果没有队伍
			// 创建队伍
			t := &model.Team{
				GameID: gameId,
				Leader: uid,
				Count:  1,
			}
			if err := dao.DB.WithContext(c.Request().Context()).Create(t).Error; err != nil {
				response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
				return
			}

			if err := dao.DB.WithContext(c.Request().Context()).Create(&model.TeamUser{
				TeamID: t.ID,
				UserID: uid,
			}).Error; err != nil {
				response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
				return
			}

			if err := dao.DB.WithContext(c.Request().Context()).Model(&model.Game{}).Where("id = ?", gameId).UpdateColumns(map[string]interface{}{
				"player_count": g.PlayerCount + 1,
				"team_count":   g.TeamCount + 1,
			}).Error; err != nil {
				response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
				return
			}
		}
	} else {
		// 查看队伍是否满员
		var t model.Team
		if err := dao.DB.WithContext(c.Request().Context()).Where("id = ?", teamId).First(&t).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
			return
		}

		if t.Count >= model.MaxTeamMember {
			response.ErrorResponse(r, model.TeamFullErrorCode, "Team full", "队伍已满")
			return
		}

		// 加入队伍
		if err := dao.DB.WithContext(c.Request().Context()).Create(&model.TeamUser{
			TeamID: teamId,
			UserID: uid,
			GameID: gameId,
		}).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
			return
		}
		if err := dao.DB.WithContext(c.Request().Context()).Model(&model.Team{}).Where("id = ?", teamId).UpdateColumns(map[string]interface{}{
			"count": g.PlayerCount + 1,
		}).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}

		if err := dao.DB.WithContext(c.Request().Context()).Model(&model.Game{}).Where("id = ?", gameId).UpdateColumns(map[string]interface{}{
			"player_count": g.PlayerCount + 1,
			"team_count":   g.TeamCount,
		}).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}
	}

	response.HTTPSuccess(r, nil)
}
