package api

import (
	"net/http"

	db "go-sip/db/sqlite"
	grpc_client "go-sip/grpc_api/c"
	. "go-sip/logger"
	"go-sip/m"
	"go-sip/model"
	sipapi "go-sip/sip"
	"go-sip/utils"
	"go-sip/zlm_api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"io"
	"strings"
)

func ZLMWebHook(c *gin.Context) {
	method := c.Param("method")
	Logger.Info("client sip ZLMWebHook", zap.String("method", method))
	switch method {
	case "on_server_started":
		// zlm启动
		ZlmServerStart(c)
	case "on_http_access":
		// http请求鉴权
	case "on_play":
		// 点播业务
		ZlmStreamOnPlay(c)
	case "on_publish":
		// 推流业务
		zlmStreamPublish(c)
	case "on_stream_none_reader":
		// 无人阅读通知 关闭流
		zlmStreamNoneReader(c)
	case "on_stream_not_found":
		// 请求播放时，流不存在时触发
		zlmStreamNotFound(c)
	case "on_record_mp4":
		//  mp4录制完成
		zlmRecordMp4(c)
	case "on_stream_changed":
		// 流注册和注销通知
		zlmStreamChanged(c)
	default:
		m.ZlmWebHookResponse(c, 0, "success")
	}

}

func ZlmServerStart(c *gin.Context) {
	m.ZlmWebHookResponse(c, 0, "zlm server start success")
}

// zlm点播业务
func ZlmStreamOnPlay(c *gin.Context) {
	m.ZlmWebHookResponse(c, 0, "on play success")
}

func zlmStreamChanged(c *gin.Context) {

	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	var req = model.ZLMStreamChangedData{}
	if err := utils.JSONDecode(data, &req); err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	if req.Regist {
		if req.APP == "rtsp" {

			Logger.Info("流注册 ", zap.Any("req", req))

			if ch, ok := grpc_client.Stream_register[req.Stream]; ok {
				ch <- struct{}{}
			}
		}
	} else {
		if req.APP == "rtmp" {
			Logger.Info("流注销 :", zap.Any("req", req))
			delete(grpc_client.Stream_register, req.Stream)
		}
	}

	m.ZlmWebHookResponse(c, 0, "success")
}

func zlmStreamPublish(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]any{
		"code":         0,
		"enable_audio": true,
		"enable_MP4":   false,
		"msg":          "success",
	})
}

func zlmRecordMp4(c *gin.Context) {
	m.ZlmWebHookResponse(c, 0, "success")
}

func zlmStreamNotFound(c *gin.Context) {

	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	var req = model.ZLMStreamNotFoundData{}
	if err := utils.JSONDecode(data, &req); err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}

	s_size := strings.Split(req.Stream, "_")

	rtp_info := zlm_api.ZlmOpenRtpServer(sipapi.Local_ZLM_Host, m.CMConfig.ZlmSecret, req.Stream, 0)
	if rtp_info.Code != 0 || rtp_info.Port == 0 {
		Logger.Error("open rtp server fail", zap.Int("code", rtp_info.Code))
		return
	}
	// 向摄像头发送信令请求推实时流到zlm
	pm := &sipapi.Streams{ChannelID: s_size[0], StreamID: s_size[0],
		ZlmIP: m.CMConfig.ZlmInnerIp, ZlmPort: rtp_info.Port, T: 0, Resolution: 1,
		Mode: 0, Ttag: db.M{}, Ftag: db.M{}, OnlyAudio: false}
	_, err = sipapi.SipPlay(pm)
	if err != nil {
		Logger.Error("向摄像头发送信令请求实时流推流到zlm失败", zap.Any("deviceId", s_size[0]), zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"code":  0,
		"close": true,
	})
}

func zlmStreamNoneReader(c *gin.Context) {
	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	var req = model.ZLMStreamNoneReaderData{}
	if err := utils.JSONDecode(data, &req); err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}

	s_size := strings.Split(req.Stream, "_")
	if len(s_size) == 3 {
		c.JSON(http.StatusOK, map[string]any{
			"code":  0,
			"close": true,
		})
	} else {
		c.JSON(http.StatusOK, map[string]any{
			"code":  0,
			"close": false,
		})
	}

}
