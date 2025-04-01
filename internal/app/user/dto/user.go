package dto

type GeneralLoginRequest struct {
	Info     string `json:"info" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type WechatRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	WechatLoginRequest
}

type GeneralLoginResponse struct {
	AccessToken          string `json:"token"`
	AccessTokenExpireIn  int64  `json:"expireIn"` // sec
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireIn int64  `json:"refreshTokenExpireIn"` // sec
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken         string `json:"token"`
	AccessTokenExpireIn int64  `json:"expireIn"` // sec
	RefreshToken        string `json:"refreshToken"`
}

type WechatLoginRequest struct {
	Code   string `json:"code" validate:"required"`
	ErrMsg string `json:"errMsg"`
}

type WechatLoginInfo struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key "`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WechatLoginResponse struct {
	AccessToken         string `json:"token"`
	AccessTokenExpireIn int64  `json:"expireIn"` // sec
	RefreshToken        string `json:"refreshToken"`
}

type MyInfoResponse struct {
	Username string `json:"username"`
}
