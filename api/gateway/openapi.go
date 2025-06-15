package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"go-sip/common"
	"go-sip/db/redis"
	"go-sip/m"
	"go-sip/model"
	"go-sip/utils"

	"io"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// @Summary     hook
// @Description zlm 启动，具体业务自行实现
func ZLMWebHook(c *gin.Context) {

	method := c.Param("method")
	if method == "" {
		c.JSON(http.StatusOK, map[string]any{
			"code": -1,
			"msg":  "method不能为空",
		})
		return
	}
	params := make(map[string]string)
	paramsArray := []string{}

	body := c.Request.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}

	switch method {
	case "on_server_started":

		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "success",
		})
		// zlm 启动
	case "on_play":
		var zlmStreamOnPlayData = model.ZlmStreamOnPlayData{}
		if err := utils.JSONDecode(bodyBytes, &zlmStreamOnPlayData); err != nil {
			m.ZlmWebHookResponse(c, -1, "body error")
			return
		}
		paramsArray = strings.Split(zlmStreamOnPlayData.Params, "&")

		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "success",
		})
	case "on_publish":
		// 推流业务
		var zlmStreamPublishData = model.ZlmStreamPublishData{}
		if err := utils.JSONDecode(bodyBytes, &zlmStreamPublishData); err != nil {
			m.ZlmWebHookResponse(c, -1, "body error")
			return
		}
		paramsArray = strings.Split(zlmStreamPublishData.Params, "&")

		c.JSON(http.StatusOK, map[string]any{
			"code":       0,
			"enableHls":  0,
			"enableMP4":  false,
			"enableRtxp": 1,
			"msg":        "success",
		})
	case "on_stream_none_reader":
		// 无人观看视频流业务

		var req = model.ZLMStreamNoneReaderData{}
		if err := utils.JSONDecode(bodyBytes, &req); err != nil {
			m.ZlmWebHookResponse(c, -1, "body error")
			return
		}

		c.JSON(http.StatusOK, map[string]any{
			"code":  0,
			"close": true,
		})

	case "on_stream_not_found":
		// 请求播放时，流不存在时触发
		var zLMStreamNotFoundData = model.ZLMStreamNotFoundData{}
		if err := utils.JSONDecode(bodyBytes, &zLMStreamNotFoundData); err != nil {
			m.ZlmWebHookResponse(c, -1, "body error")
			return
		}
		paramsArray = strings.Split(zLMStreamNotFoundData.Params, "&")

		c.JSON(http.StatusOK, map[string]any{
			"code":  0,
			"close": true,
		})
	case "on_record_mp4":
		//  mp4 录制完成
		zlmRecordMp4(c)
		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "success",
		})
	case "on_stream_changed":
		// 流注册和注销通知
		var zLMStreamChangedData = model.ZLMStreamChangedData{}
		if err := utils.JSONDecode(bodyBytes, &zLMStreamChangedData); err != nil {
			m.ZlmWebHookResponse(c, -1, "body error")
			return
		}
		paramsArray = strings.Split(zLMStreamChangedData.Params, "&")
		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "success",
		})
	default:
		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "success",
		})
		return
	}

	if len(paramsArray) != 0 {
		for _, param := range paramsArray {
			if param == "" {
				continue
			}
			tokens := strings.Split(param, "=")
			if len(tokens) != 2 {
				return
			}
			params[tokens[0]] = tokens[1]
		}
	}

	server_url := ""
	if ipc, ok := params["stream"]; ok {

		stream_id_split := strings.Split(ipc, "_")

		ipc_id := stream_id_split[0]

		sip_server, err := redis.HGet(c.Copy(), redis.IPC_SIPSERVER, ipc_id)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": -1,
				"msg":  "ipc 未注册",
			})
			return
		}

		server_url, err = redis.HGet(c.Copy(), redis.SIPSERVER_Node, sip_server)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": -1,
				"msg":  "sip server 未找到",
			})
			return
		}

	} else {
		sip_server_map, err := redis.HGetAll(c.Copy(), redis.SIPSERVER_Node)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": -1,
				"msg":  "not found sip server",
			})
			return
		}

		_, sip_server_url, err := utils.SelectRandomMapValue(sip_server_map)
		if err != nil {
			c.JSON(http.StatusOK, map[string]any{
				"code": -1,
				"msg":  "not found sip server",
			})
			return
		}

		server_url = sip_server_url

	}

	// // 使用sip_url调用sip服务接口
	full_url := fmt.Sprintf("%s%s", server_url, common.ZLMWebHookBaseURL+"/"+method)

	// 调用sip接口
	req, err := http.NewRequest("POST", full_url, bytes.NewReader(bodyBytes))
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

	switch method {
	case "on_http_access":
		c.JSON(http.StatusOK, map[string]any{
			"code":   0,
			"second": 86400})
	case "on_publish":
		// 推流鉴权
		c.JSON(http.StatusOK, map[string]any{
			"code":       0,
			"enableHls":  1,
			"enableMP4":  false,
			"enableRtxp": 1,
			"msg":        "success",
		})
	case "on_stream_not_found":
		c.JSON(http.StatusOK, map[string]any{
			"code":  0,
			"close": true,
		})
	case "on_stream_none_reader":
		c.JSON(http.StatusOK, map[string]any{
			"code":  0,
			"close": true,
		})
	default:
		c.JSON(http.StatusOK, map[string]any{
			"code": 0,
			"msg":  "success",
		})
	}

}

func RegisterSipServer(c *gin.Context) {

	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "body错误")
		return
	}
	t := model.SipServerInfo{}
	err = utils.JSONDecode(data, &t)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "反序列化失败，body错误")
		return
	}

	url := fmt.Sprintf("http://%s:%s", c.RemoteIP(), t.Port)

	redis.HSet(c.Copy(), redis.SIPSERVER_Node, t.ServerID, url)

	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "success",
	})
}

func GetSipServer(c *gin.Context) {

}

// @Summary     ipc回放视频倍速
// @Description 用来设置通道设备回放视频的播放速度，注意控制速度范围 0.25- 4.0，超过范围会导致回放失败。
// @Tags        playbackSpeed
// @Router      /open/ipc/playbackSpeed [get]
// @Param       ipc_id    path     string true "通道id"
// @Param       speed query    float    true 	"回放速度，范围0.25-4.0，默认1.0"
// @Success     0
func PlaybackSpeed(c *gin.Context) {
	stream_id := c.Query("stream")

	params := url.Values{}
	params.Add("stream", stream_id)
	params.Add("speed", c.Query("speed"))

	stream_id_split := strings.Split(stream_id, "_")

	ipc_id := stream_id_split[0]

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
	full_url := fmt.Sprintf("%s%s?%s", url, common.PlaybackURL, params.Encode())

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

	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "success",
	})

	// 获取sip_id
}

// @Summary     ipc回放视频拖动播放
// @Description 用来设置通道设备回放视频的拖动播放位置，注意控制拖动时间范围，超过范围会导致回放失败。
// @Tags        PlaybackSeek
// @Router      /open/ipc/playbackSeek [get]
// @Param       ipc_id    path     string true "通道id"
// @Param       seek query    int    true  "拖动时间，单位秒"
// @Success     0
func PlaybackSeek(c *gin.Context) {
	stream_id := c.Query("stream")

	params := url.Values{}
	params.Add("stream", stream_id)
	params.Add("seek", c.Query("seek"))

	stream_id_split := strings.Split(stream_id, "_")

	ipc_id := stream_id_split[0]

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
	full_url := fmt.Sprintf("%s%s?%s", url, common.SeekURL, params.Encode())

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

	c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "success",
	})

}

// @Summary     ipc回放视频停止
// @Description 用来停止通道设备的回放视频
// @Tags        PlaybackStop
// @Router      /open/ipc/playbackStop [get]
// @Param       ipc_id    path     string true "通道id"
// @Success     0
func PlaybackStop(c *gin.Context) {
	stream_id := c.Query("stream")

	stream_id_split := strings.Split(stream_id, "_")

	ipc_id := stream_id_split[0]

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
	full_url := fmt.Sprintf("%s%s?stream=%s", url, common.StopURL, stream_id)

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

// @Summary     ipc回放视频暂停
// @Description 用来暂停通道设备的回放视频
// @Tags        PlaybackPause
// @Router      /open/ipc/playbackPause [get]
// @Param       ipc_id    path     string true "通道id"
// @Success     0
func PlaybackPause(c *gin.Context) {

	stream_id := c.Query("stream")

	stream_id_split := strings.Split(stream_id, "_")

	ipc_id := stream_id_split[0]

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
	full_url := fmt.Sprintf("%s%s?stream=%s", url, common.PauseURL, stream_id)

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

// @Summary     ipc回放视频恢复
// @Description 用来恢复通道设备的回放视频
// @Tags        PlaybackResume
// @Router      /open/ipc/playbackResume [get]
// @Param       ipc_id    path     string true "通道id"
// @Success     0
func PlaybackResume(c *gin.Context) {

	stream_id := c.Query("stream")

	stream_id_split := strings.Split(stream_id, "_")

	ipc_id := stream_id_split[0]

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
	full_url := fmt.Sprintf("%s%s?stream=%s", url, common.ResumeURL, stream_id)

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

// @Summary     ipc回放视频时间列表
// @Description 用来获取通道设备存储的可回放时间段列表，注意控制时间跨度，跨度越大，数据量越多，返回越慢，甚至会超时（最多10s）。
// @Tags        RecordsList
// @Router      /open/ipc/records [get]
// @Param       ipc_id    path     string true "通道id"
// @Param       start query    int    true "开始时间，时间戳"
// @Param       end   query    int    true "结束时间，时间戳"
// @Success     0     {object} sipapi.Records
// @Failure     1000  {object} string
// @Failure     1001  {object} string
// @Failure     1002  {object} string
// @Failure     1003  {object} string
func RecordsList(c *gin.Context) {
	ipc_id := c.Query("ipc_id")

	params := url.Values{}
	params.Add("ipc_id", ipc_id)
	params.Add("start", c.Query("start"))
	params.Add("end", c.Query("end"))

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
	full_url := fmt.Sprintf("%s%s?%s", url, common.RecordsListURL, params.Encode())
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

	data := response.Data
	// json反序列化
	var dataContent model.DataContent
	err = json.Unmarshal([]byte(data.(string)), &dataContent)
	if err != nil {
		model.JsonResponseSysERR(c, "json反序列化失败")
		return
	}
	list := dataContent.List
	if len(list) == 0 {
		model.JsonResponseSysERR(c, nil)
		return
	}

	model.JsonResponseSucc(c, list[0].Items)

}
