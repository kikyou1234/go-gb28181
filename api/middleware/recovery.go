package middleware

import (
	"go-sip/m"

	"github.com/gin-gonic/gin"
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			m.JsonResponse(c, c.Err().Error(), "服务器错误，请联系管理员")
		}
	}()
	c.Next()
}
