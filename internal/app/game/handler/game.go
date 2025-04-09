package handler

import (
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/auth"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/logx"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/flamego"
)

func HandleGameCreate(c flamego.Context, r flamego.Render, auth auth.Info) {
	var uid string
	if uid = auth.Uid; uid == "" {
		response.ErrorResponse(r, model.UnauthorizedErrorCode, "Unauthorized", "未授权")
		return
	}

	// 一个人只能加入一个游戏
	var tu []model.TeamUser
	if err := dao.DB.WithContext(c.Request().Context()).Where("user_id = ?", uid).Find(&tu).Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}

	if len(tu) > 0 {
		response.ErrorResponse(r, model.GameAlreadyJoinErrorCode, "Game already joined", "已经加入游戏")
		return
	}

	tx := dao.DB.WithContext(c.Request().Context()).Begin()
	defer tx.Rollback()

	var g = &model.Game{
		Sponsor: uid,
	}
	if err := tx.Model(&model.Game{}).Create(g).Error; err != nil {
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	t := &model.Team{
		GameID: g.ID,
		Leader: uid,
	}
	if err := tx.Model(&model.Team{}).Create(t).Error; err != nil {
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	if err := tx.Model(&model.TeamUser{}).Create(&model.TeamUser{
		TeamID: t.ID,
		UserID: uid,
	}).Error; err != nil {
		response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.ErrorResponse(r, model.DatabaseTransactionErrorCode, err)
		return
	}

	response.HTTPSuccess(r, g)
}

func HandleGameJoin(c flamego.Context, r flamego.Render, req dto.GameJoinRequest, auth auth.Info) {
	var uid string
	if uid = auth.Uid; uid == "" {
		response.ErrorResponse(r, model.UnauthorizedErrorCode, "Unauthorized", "未授权")
		return
	}

	code := req.Code

	// 一个人只能加入一个游戏
	var tu []model.TeamUser
	if err := dao.DB.WithContext(c.Request().Context()).Where("user_id = ?", uid).Find(&tu).Error; err != nil {
		logx.NameSpace("game").Error(err)
		response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		return
	}
	var g model.Game
	if len(tu) > 0 {
		if err := dao.DB.WithContext(c.Request().Context()).Where("code = ?", code).First(&g).Error; err != nil {
			logx.NameSpace("game").Error(err)
			response.ErrorResponse(r, model.DatabaseFirstErrorCode, err)
			return
		}
		// 这个游戏满员了
		if g.PlayerCount >= model.MaxPlayerCount {
			response.ErrorResponse(r, model.GameFullErrorCode, "Game full", "游戏已满员")
			return
		}

		// 已经加入了这个游戏
		if tu[0].GameID == g.ID {
			response.HTTPSuccess(r, nil)
			return
		}

		// 已经加入了其他游戏
		if tu[0].GameID != g.ID {
			response.ErrorResponse(r, model.GameAlreadyJoinErrorCode, "Game already joined", "已经加入游戏")
			return
		}
	} else {
		// 如果没有队伍就直接加入这个游戏; 如果有队伍就加入这个游戏的队伍
		tx := dao.DB.WithContext(c.Request().Context()).Begin()
		defer tx.Rollback()

		var t = &model.Team{
			GameID: g.ID,
			Leader: uid,
		}
		if err := tx.Create(t).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
			return
		}

		if err := tx.Create(&model.TeamUser{
			TeamID: t.ID,
			UserID: uid,
			GameID: g.ID,
		}).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
			return
		}

		if err := tx.Model(&model.Game{}).Where("id = ?", g.ID).Update("player_count", g.PlayerCount+1).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(r, model.DatabaseTransactionErrorCode, err)
			return
		}
	}

	// 这里需要返回游戏的详细信息
	response.HTTPSuccess(r, nil)
}

func HandleGameExit(c flamego.Context, r flamego.Render, auth auth.Info) {
}
