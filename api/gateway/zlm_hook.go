package gateway

import (
	"go-sip/m"
	"go-sip/model"
	"go-sip/utils"
	"io"

	"github.com/gin-gonic/gin"
)

func zlmRecordMp4(c *gin.Context) {
	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	var req = model.ZLMRecordMp4Data{}
	if err := utils.JSONDecode(data, &req); err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}

	m.ZlmWebHookResponse(c, 0, "success")
}
