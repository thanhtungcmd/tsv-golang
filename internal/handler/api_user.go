package handler

import (
	"tsv-golang/internal/dto"
	"tsv-golang/internal/persistence"

	"github.com/gin-gonic/gin"
)

type ApiUserHandler struct {
	repo persistence.Repositories
}

func ApiUserHandlerInit(repo persistence.Repositories) *ApiUserHandler {
	return &ApiUserHandler{
		repo: repo,
	}
}

func (h *ApiUserHandler) GetList(c *gin.Context) {
	params := &dto.GetListUserRequest{}
	_ = c.ShouldBindQuery(&params)

	LogInfo(c, "UserGetListRequest:", params)

	data := h.repo.User.GetList(params)

	LogInfo(c, "UserGetListResponse:", len(data))

	ResponseSuccess(c, data)
	return
}
