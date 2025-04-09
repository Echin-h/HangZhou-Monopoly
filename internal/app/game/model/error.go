package model

const (
	// UnauthorizedCode 权限报错
	UnauthorizedErrorCode = 40301
	ParamErrorCode        = 40302

	// DatabaseCode 数据库报错
	DatabaseCreateErrorCode = 50000 + iota
	DatabaseFirstErrorCode
	DatabaseFindErrorCode
	DatabaseUpdateErrorCode
	DatabaseDeleteErrorCode
	DatabaseTransactionErrorCode

	// GameCode 游戏相关报错
	GameAlreadyJoinErrorCode = 50100 + iota
	GameFullErrorCode

	// TeamCode 队伍相关报错
	TeamAlreadyJoinErrorCode = 50200 + iota
	TeamFullErrorCode
)
