package gateway

import (
	"encoding/json"
	"fmt"
	"go-sip/common"
	"go-sip/db/redis"
	"go-sip/m"
	"go-sip/model"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func DeviceControl(c *gin.Context) {
	ipc_id := c.Query("ipc_id")

	params := url.Values{}
	params.Add("ipc_id", ipc_id)
	params.Add("leftRight", c.Query("leftRight"))
	params.Add("upDown", c.Query("upDown"))
	params.Add("inOut", c.Query("inOut"))
	params.Add("moveSpeed", c.Query("moveSpeed"))
	params.Add("zoomSpeed", c.Query("zoomSpeed"))

	sip_server, err := redis.HGet(c.Copy(), redis.IPC_SIPSERVER, ipc_id)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"code": -1,
			"msg":  "ipc 未注册",
		})
		return
	}

	url, err := redis.HGet(c.Copy(), redis.SIPSERVER_Node, sip_server)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"code": -1,
			"msg":  "sip server 未找到",
		})
		return
	}

	// // 使用sip_url调用sip服务接口
	full_url := fmt.Sprintf("%s%s?%s", url, common.DeviceControl, params.Encode())

	// 调用sip接口
	req, err := http.NewRequest("GET", full_url, nil)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{
			"code": -1,
			"msg":  "json marshal error",
		})
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		model.JsonResponseSysERR(c, "调用sip hook接口失败")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		model.JsonResponseSysERR(c, "调用sip hook接口失败")
	}
	var response m.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		model.JsonResponseSysERR(c, "json反序列化失败")
	}

	c.JSON(http.StatusOK, map[string]any{
		"code": response.Code,
		"msg":  response.Data,
	})
}
