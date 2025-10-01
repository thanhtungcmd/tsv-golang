package handler

import (
	"encoding/json"
	"tsv-golang/internal/dto"
	"tsv-golang/internal/persistence"
	"tsv-golang/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
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

	paramMap := make(map[string]interface{})
	paramMap["TrackID"], _ = c.Get("TrackID")
	_ = mapstructure.Decode(params, &paramMap)
	paramBytes, _ := json.Marshal(paramMap)
	log.Info("User Request Params:", string(paramBytes))

	data := h.repo.User.GetList(params)
	ResponseSuccess(c, data)
	return
}
