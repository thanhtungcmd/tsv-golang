package dto

type GetListUserRequest struct {
	Limit int64  `json:"limit" form:"limit"`
	Page  int64  `json:"page" form:"page"`
	Id    string `json:"id" form:"id"`
}
