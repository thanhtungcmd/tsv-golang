package route

import (
	"tsv-golang/internal/handler"
	"tsv-golang/internal/persistence"

	"github.com/gin-gonic/gin"
)

func HandleApiV1(route *gin.RouterGroup, repo persistence.Repositories) {
	paygateHandler := handler.ApiUserHandlerInit(repo)
	route.GET("users/get-list", paygateHandler.GetList)
}
