package handler

import (
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
	
	ResponseSuccess(c, "123")
	return
}
