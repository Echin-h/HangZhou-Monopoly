package dto

type TeamCreateRequest struct{}

type TeamJoinRequest struct {
	Code int `json:"code" binding:"required"`
}
