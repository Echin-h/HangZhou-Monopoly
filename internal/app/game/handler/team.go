package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/auth"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/flamego"
	"gorm.io/gorm"
)

// 当前游戏里的人加入指定的队伍，若指定的teamid为空或者不存在，则创建一个队伍
func HandleTeamJoin(c flamego.Context, r flamego.Render, auth auth.Info) {
	// 1. 用户认证检查
	uid := auth.Uid
	if uid == "" {
		response.ErrorResponse(r, model.UnauthorizedErrorCode, "Unauthorized", "未授权")
		return
	}

	// 2. 参数校验
	gameId := c.Query("game_id")
	if gameId == "" {
		response.ErrorResponse(r, model.ParamErrorCode, "Game ID empty", "游戏ID不能为空")
		return
	}
	teamId := c.Query("team_id")

	// 3. 检查游戏状态
	var g model.Game
	if err := dao.DB.WithContext(c.Request().Context()).
		Where("id = ?", gameId).
		First(&g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(r, model.GameNotFoundErrorCode, "Game not found", "游戏不存在")
		} else {
			response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		}
		return
	}

	if g.PlayerCount >= model.MaxPlayerCount {
		response.ErrorResponse(r, model.GameFullErrorCode, "Game full", "游戏人数已满")
		return
	}

	// 4. 检查用户是否已加入游戏（必须有TeamUser记录才能操作）
	var existingTeamUser model.TeamUser
	if err := dao.DB.WithContext(c.Request().Context()).
		Where("user_id = ? AND game_id = ?", uid, gameId).
		First(&existingTeamUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(r, model.TeamAlreadyJoinErrorCode,
				"Must join game first", "请先加入游戏")
		} else {
			response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
		}
		return
	}

	tx := dao.DB.WithContext(c.Request().Context()).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var finalTeamID string
	var isNewTeam bool

	// 5. 处理队伍加入/创建逻辑
	if teamId != "" {
		// 5.1 尝试加入指定队伍
		var targetTeam model.Team
		if err := tx.Where("id = ? AND game_id = ?", teamId, gameId).
			First(&targetTeam).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(r, model.TeamNotFoundErrorCode, "Team not found", "队伍不存在")
			} else {
				response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
			}
			return
		}

		if targetTeam.Count >= model.MaxTeamMember {
			response.ErrorResponse(r, model.TeamFullErrorCode, "Team full", "队伍已满")
			return
		}
		finalTeamID = targetTeam.ID
		isNewTeam = false
	} else {
		// 5.2 创建新队伍（检查队伍数量限制）
		// 检查当前游戏的队伍数量
		var teamCount int64
		if err := tx.Model(&model.Team{}).
			Where("game_id = ?", gameId).
			Count(&teamCount).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseFindErrorCode, err)
			return
		}

		// 验证队伍数量上限（最多4个队伍）
		if teamCount >= 4 {
			response.ErrorResponse(r, model.TeamLimitErrorCode,
				"Team limit reached", "队伍数量已达上限（最多4个）")
			return
		}

		// 生成队伍ID格式：游戏ID-队伍序号（1-4）
		teamNumber := teamCount + 1
		teamID := fmt.Sprintf("%s-%d", gameId, teamNumber)

		// 创建新队伍记录
		newTeam := &model.Team{
			ID:        teamID, // 格式：游戏ID-1|2|3|4
			GameID:    gameId,
			Name:      fmt.Sprintf("队伍%d", teamNumber), // 默认队伍名称
			Balance:   0,                               // 初始金额
			Leader:    uid,                             // 队长为当前用户
			Count:     1,                               // 初始人数
			CreatedAt: time.Now(),
		}

		// 将新队伍插入teams表
		if err := tx.Create(newTeam).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseCreateErrorCode, err)
			return
		}
		finalTeamID = newTeam.ID
		isNewTeam = true

		// 更新游戏的队伍计数
		if err := tx.Model(&model.Game{}).
			Where("id = ?", gameId).
			Update("team_count", teamNumber).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}
	}

	// 6. 更新用户队伍关系（而不是创建新记录）
	if err := tx.Model(&model.TeamUser{}).
		Where("user_id = ? AND game_id = ?", uid, gameId).
		Update("team_id", finalTeamID).Error; err != nil {
		response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
		return
	}

	// 7. 更新队伍人数（只有加入已有队伍时才增加人数）
	if !isNewTeam {
		if err := tx.Model(&model.Team{}).
			Where("id = ?", finalTeamID).
			Update("count", gorm.Expr("count + 1")).Error; err != nil {
			response.ErrorResponse(r, model.DatabaseUpdateErrorCode, err)
			return
		}
	}

	// 8. 提交事务
	if err := tx.Commit().Error; err != nil {
		response.ErrorResponse(r, model.DatabaseTransactionErrorCode, err)
		return
	}

	// 9. 返回成功响应
	response.HTTPSuccess(r, map[string]interface{}{
		"game_id": gameId,
		"team_id": finalTeamID,
		"action":  map[bool]string{true: "created", false: "joined"}[isNewTeam],
	})
}
