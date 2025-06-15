package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	StatusSucc   = 200
	StatusSysERR = 500
)

const (
	CodeSucc   = "success"
	CodeSysERR = "fail"
)

var HTTPCode = map[string]int{
	CodeSucc: http.StatusOK,
}

type ApiResult struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
	Result any    `json:"result"`
}

func JsonResponse(c *gin.Context, code string, status int, data any) {
	switch d := data.(type) {
	case error:
		data = d.Error()
	}
	c.JSON(HTTPCode[code], ApiResult{Code: code, Status: status, Result: data})
}

func JsonResponseSucc(c *gin.Context, data any) {
	JsonResponse(c, CodeSucc, StatusSucc, data)
}

func JsonResponseSysERR(c *gin.Context, data any) {
	JsonResponse(c, CodeSysERR, StatusSysERR, data)
}
