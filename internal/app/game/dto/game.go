package dto

type ExamplePost struct {
	Msg string `json:"msg"`
}
type GameJoinRequest struct {
	Code int `json:"code" binding:"required"`
}
