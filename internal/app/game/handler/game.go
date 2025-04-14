package handler

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/auth"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/logx"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/flamego"
	"github.com/oklog/ulid/v2"
)

func HandleGameCreate(c flamego.Context, r flamego.Render, auth auth.Info) {
	uid := auth.Uid
	if uid == "" {
		response.ErrorResponse(r, model.UnauthorizedErrorCode, "Unauthorized", "未授权")
		return
	}

	// 检查是否已加入游戏
	var tu []model.TeamUser
	if err := dao.DB.WithContext(c.Request().Context()).
		Where("user_id = ?", uid).
		Find(&tu).Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}

	if len(tu) > 0 {
		response.ErrorResponse(r, model.GameAlreadyJoinErrorCode,
			"Game already joined", "已经加入游戏")
		return
	}

	tx := dao.DB.WithContext(c.Request().Context()).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建游戏
	game := &model.Game{
		Sponsor: uid,
	}
	if err := tx.Create(game).Error; err != nil {
		tx.Rollback()
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	// 创建团队
	team := &model.Team{
		ID:     ulid.Make().String(), // 生成唯一ID
		GameID: game.ID,              // 关联游戏ID
		Leader: uid,
		Count:  1,
	}
	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	// 创建团队成员（关键修复：添加GameID）
	teamUser := &model.TeamUser{
		TeamID: team.ID,
		UserID: uid,
		GameID: game.ID, // 这里设置GameID
	}
	if err := tx.Create(teamUser).Error; err != nil {
		tx.Rollback()
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	// 更新游戏人数
	if err := tx.Model(&model.Game{}).
		Where("id = ?", game.ID).
		Update("player_count", 1).Error; err != nil {
		tx.Rollback()
		response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorResponse(r, model.DatabaseTransactionErrorCode, err)
		return
	}

	response.HTTPSuccess(r, game)
}
func HandleGameJoin(c flamego.Context, r flamego.Render, req dto.GameJoinRequest, auth auth.Info) {
	uid := auth.Uid
	if uid == "" {
		response.ErrorResponse(r, model.UnauthorizedErrorCode, "Unauthorized", "未授权")
		return
	}

	code := req.Code

	// 1. 先查询目标游戏
	var targetGame model.Game
	if err := dao.DB.WithContext(c.Request().Context()).
		Where("code = ?", code).
		First(&targetGame).Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseFirstErrorCode, err)
		return
	}

	// 2. 检查玩家是否已加入其他游戏
	var existingTeamUsers []model.TeamUser
	if err := dao.DB.WithContext(c.Request().Context()).
		Where("user_id = ?", uid).
		Find(&existingTeamUsers).Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}

	// 3. 已加入其他游戏的检查
	if len(existingTeamUsers) > 0 {
		if existingTeamUsers[0].GameID == targetGame.ID {
			// 已加入当前游戏
			response.HTTPSuccess(r, nil)
			return
		}
		// 已加入其他游戏
		response.ErrorResponse(r, model.GameAlreadyJoinErrorCode,
			"Game already joined", "已经加入游戏")
		return
	}

	// 4. 检查游戏是否满员
	if targetGame.PlayerCount >= model.MaxPlayerCount {
		response.ErrorResponse(r, model.GameFullErrorCode,
			"Game full", "游戏已满员")
		return
	}

	// 5. 检查游戏是否有队伍
	var existingTeams []model.Team
	if err := dao.DB.WithContext(c.Request().Context()).
		Where("game_id = ?", targetGame.ID).
		Find(&existingTeams).Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}

	tx := dao.DB.WithContext(c.Request().Context()).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var teamID string

	if len(existingTeams) > 0 {
		// 有队伍，加入第一个未满的队伍
		var joinTeam *model.Team
		for _, team := range existingTeams {
			if team.Count < model.MaxTeamMember {
				joinTeam = &team
				break
			}
		}

		if joinTeam == nil {
			tx.Rollback()
			response.ErrorResponse(r, model.TeamFullErrorCode,
				"All teams are full", "所有队伍都已满")
			return
		}

		teamID = joinTeam.ID

		// 更新队伍人数
		if err := tx.Model(&model.Team{}).
			Where("id = ?", teamID).
			Update("count", joinTeam.Count+1).Error; err != nil {
			tx.Rollback()
			logx.NameSpace("game").Error(err)
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}
	} else {
		// 没有队伍，创建新队伍
		teamID = ulid.Make().String()
		newTeam := &model.Team{
			ID:     teamID,
			GameID: targetGame.ID,
			Leader: uid,
			Count:  1,
		}
		if err := tx.Create(newTeam).Error; err != nil {
			tx.Rollback()
			logx.NameSpace("game").Error(err)
			response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
			return
		}

		// 更新游戏队伍数
		if err := tx.Model(&model.Game{}).
			Where("id = ?", targetGame.ID).
			Update("team_count", targetGame.TeamCount+1).Error; err != nil {
			tx.Rollback()
			logx.NameSpace("game").Error(err)
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}
	}

	// 添加团队成员
	if err := tx.Create(&model.TeamUser{
		TeamID: teamID,
		UserID: uid,
		GameID: targetGame.ID,
	}).Error; err != nil {
		tx.Rollback()
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	// 更新游戏人数
	if err := tx.Model(&model.Game{}).
		Where("id = ?", targetGame.ID).
		Update("player_count", targetGame.PlayerCount+1).Error; err != nil {
		tx.Rollback()
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseTransactionErrorCode, err)
		return
	}

	response.HTTPSuccess(r, map[string]interface{}{
		"game_id": targetGame.ID,
		"team_id": teamID,
	})
}

func HandleGameExit(c flamego.Context, r flamego.Render, auth auth.Info) {
}
