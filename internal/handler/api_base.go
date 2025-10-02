package handler

import (
	"encoding/json"
	"net/http"
	"tsv-golang/pkg/error_common"
	"tsv-golang/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type ResponseSuccessStruct struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}
type ResponseFailStruct struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	Errors    interface{} `json:"errors,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}
type ResponseStruct struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	HttpCode  int         `json:"httpCode,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	any       any
}

func ResponseSuccess(context *gin.Context, data interface{}) {
	res := ResponseSuccessStruct{}
	res.ErrorCode = error_common.SUCCESS
	res.Message = error_common.Error(error_common.SUCCESS).GetText()
	res.Data = data
	context.JSON(http.StatusOK, res)
}

func ResponseSuccessMessage(context *gin.Context, message string, data interface{}) {
	res := ResponseSuccessStruct{}
	res.ErrorCode = error_common.SUCCESS
	res.Message = message
	res.Data = data
	context.JSON(http.StatusOK, res)
}

func ResponseFail(context *gin.Context, errorCode int, errors ...interface{}) {
	res := ResponseFailStruct{}
	res.ErrorCode = errorCode
	res.Message = error_common.Error(errorCode).GetText()
	if errors != nil {
		res.Errors = errors
	}
	context.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func ResponseFailWithData(context *gin.Context, errorCode int, data interface{}) {
	res := ResponseFailStruct{}
	res.ErrorCode = errorCode
	res.Message = error_common.Error(errorCode).GetText()
	res.Data = data
	context.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func Response(context *gin.Context, res any, httpCode int) {
	context.JSON(httpCode, res)
}

func LogInfo(c *gin.Context, content string, params interface{}) {
	paramMap := make(map[string]interface{})
	paramMap["TrackID"], _ = c.Get("TrackID")
	var paramTemp interface{}
	_ = mapstructure.Decode(params, &paramTemp)
	paramMap["Data"] = paramTemp
	paramBytes, _ := json.Marshal(paramMap)
	log.Info(content, string(paramBytes))
}
