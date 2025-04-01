package handler

import (
	"time"

	"github.com/Echin-h/HangZhou-Monopoly/config"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dao"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/dto"
	"github.com/Echin-h/HangZhou-Monopoly/internal/app/user/model"
	"github.com/Echin-h/HangZhou-Monopoly/internal/core/auth"
	"github.com/Echin-h/HangZhou-Monopoly/internal/middleware/response"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/levigross/grequests"
)

func HandleWechatLogin(r flamego.Render, c flamego.Context, req dto.WechatLoginRequest, errs binding.Errors) {
	if errs != nil {
		response.ServiceErr(r, errs)
		return
	}
	code := req.Code
	url := "https://api.weixin.qq.com/sns/jscode2session"
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Params: map[string]string{
			"appid":      config.GetConfig().WxLogin.AppId,
			"secret":     config.GetConfig().WxLogin.AppSecret,
			"grant_type": "authorization_code",
			"js_code":    code,
		},
		Context: c.Request().Context(),
	})
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	if !resp.Ok {
		response.ServiceErr(r, resp.Error)
		return
	}

	var info dto.WechatLoginInfo
	err = resp.JSON(&info)
	if err != nil {
		response.ServiceErr(r, err)
		return
	}

	if info.ErrCode != 0 {
		switch info.ErrCode {
		case 40029:
			response.ServiceErr(r, "invalid code:"+info.ErrMsg)
			return
		case 45011:
			response.ServiceErr(r, "frequent request:"+info.ErrMsg)
			return
		default:
			response.ServiceErr(r, "unknown error:"+info.ErrMsg)
			return
		}
	}

	//查找User表，确认是否已经注册
	var user model.User
	err = dao.DB.WithContext(c.Request().Context()).Where("open_id = ?", info.OpenId).Limit(1).Find(&user).Error
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	if user.OpenId != "" {
		//已经注册，返回token
		var token string
		token, err = auth.GenToken(auth.Info{
			Uid:        user.OpenId,
			SessionKey: info.SessionKey,
		})
		if err != nil {
			response.ServiceErr(r, err)
			return
		}
		var refreshToken string
		refreshToken, err = auth.GenToken(auth.Info{
			Uid:            user.OpenId,
			SessionKey:     info.SessionKey,
			IsRefreshToken: true,
		}, auth.RefreshTokenExpireIn)
		if err != nil {
			response.ServiceErr(r, err)
			return
		}
		response.HTTPSuccess(r, dto.WechatLoginResponse{
			AccessToken:         token,
			AccessTokenExpireIn: int64(auth.AccessTokenExpireIn / time.Second),
			RefreshToken:        refreshToken,
		})
		return
	}
}

func HandleTest(r flamego.Render, c flamego.Context) {
	var user model.User
	err := dao.DB.Model(&model.User{}).WithContext(c.Request().Context()).Where("open_id = 'aaa'").First(&user).Error
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	var token string
	token, err = auth.GenToken(auth.Info{
		Uid: user.OpenId,
	})
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	var refreshToken string
	refreshToken, err = auth.GenToken(auth.Info{
		Uid:            user.OpenId,
		IsRefreshToken: true,
	}, auth.RefreshTokenExpireIn)
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	response.HTTPSuccess(r, dto.WechatLoginResponse{
		AccessToken:         token,
		AccessTokenExpireIn: int64(auth.AccessTokenExpireIn / time.Second),
		RefreshToken:        refreshToken,
	})
	return

}

func HandleRefreshToken(r flamego.Render, req dto.RefreshTokenRequest) {
	entity, err := auth.ParseToken(req.RefreshToken)
	if err != nil || !entity.Info.IsRefreshToken {
		response.UnAuthorization(r)
		return
	}

	token, err := auth.GenToken(auth.Info{Uid: entity.Info.Uid})
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	response.HTTPSuccess(r, dto.RefreshTokenResponse{
		AccessToken:         token,
		AccessTokenExpireIn: int64(auth.AccessTokenExpireIn / time.Second),
		RefreshToken:        req.RefreshToken,
	})
}

func HandleRegister(r flamego.Render, c flamego.Context, req dto.WechatRegisterRequest, errs binding.Errors) {
	if errs != nil {
		response.ServiceErr(r, errs)
		return
	}
	code := req.Code
	url := "https://api.weixin.qq.com/sns/jscode2session"
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Params: map[string]string{
			"appid":      config.GetConfig().WxLogin.AppId,
			"secret":     config.GetConfig().WxLogin.AppSecret,
			"grant_type": "authorization_code",
			"js_code":    code,
		},
		Context: c.Request().Context(),
	})
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	if !resp.Ok {
		response.ServiceErr(r, resp.Error)
		return
	}

	var info dto.WechatLoginInfo
	err = resp.JSON(&info)
	if err != nil {
		response.ServiceErr(r, err)
		return
	}

	if info.ErrCode != 0 {
		switch info.ErrCode {
		case 40029:
			response.ServiceErr(r, "invalid code:"+info.ErrMsg)
			return
		case 45011:
			response.ServiceErr(r, "frequent request:"+info.ErrMsg)
			return
		default:
			response.ServiceErr(r, "unknown error:"+info.ErrMsg)
			return
		}
	}

	user := model.User{
		Username: req.Username,
		OpenId:   info.OpenId,
	}
	err = dao.DB.WithContext(c.Request().Context()).Create(&user).Error
	if err != nil {
		response.ServiceErr(r, err)
		return
	}

	var token string
	token, err = auth.GenToken(auth.Info{
		Uid:        info.OpenId,
		SessionKey: info.SessionKey,
	})
	if err != nil {
		response.ServiceErr(r, err)
		return
	}

	var refreshToken string
	refreshToken, err = auth.GenToken(auth.Info{
		Uid:            info.OpenId,
		SessionKey:     info.SessionKey,
		IsRefreshToken: true,
	}, auth.RefreshTokenExpireIn)
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	response.HTTPSuccess(r, dto.WechatLoginResponse{
		AccessToken:         token,
		AccessTokenExpireIn: int64(auth.AccessTokenExpireIn / time.Second),
		RefreshToken:        refreshToken,
	})
}

func HandleGetMyInfo(r flamego.Render, c flamego.Context, info auth.Info) {
	var user model.User
	err := dao.DB.WithContext(c.Request().Context()).Where("open_id = ?", info.Uid).Limit(1).Find(&user).Error
	if err != nil {
		response.ServiceErr(r, err)
		return
	}
	response.HTTPSuccess(r, dto.MyInfoResponse{
		Username: user.Username,
	})
}

func HandlePing(r flamego.Render) {
	response.HTTPSuccess(r, "pong")
}
