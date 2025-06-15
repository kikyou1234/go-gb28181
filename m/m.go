package m

import (
	"go-sip/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	StatusSucc      = "0"
	StatusAuthERR   = "1000"
	StatusDBERR     = "1001"
	StatusParamsERR = "1002"
	StatusSysERR    = "1003"

	StreamTypePull = "pull"
	StreamTypePush = "push"
)

const (
	Ping          = "ping"
	Play          = "play"      // 点播
	PlayBack      = "play_back" // 回放点播
	StopPlay      = "stop_play"
	PausePlay     = "pause_play"
	ResumePlay    = "resume_play"
	SpeedPlay     = "speed_play"  // 倍速
	SeekPlay      = "seek_play"   //进度条
	RecordList    = "record_list" // 获取回放历史记录
	Broadcast     = "broadcast"   // 语音广播
	PlayIPCAudio  = "play_Ipc_audio"
	DeviceControl = "ipc_device_control"

	PlayAudio  = "play_audio"
	PushAudio  = "push_audio"
	SetVolume  = "set_volume"
	CloseAudio = "close_audio"
)

const (
	MsgID_Ping         = 0
	MsgID_Play         = 1
	MsgID_PlayBack     = 2
	MsgID_StopPlay     = 3
	MsgID_PausePlay    = 4
	MsgID_ResumePlay   = 5
	MsgID_SpeedPlay    = 6
	MsgID_SeekPlay     = 7
	MsgID_RecordList   = 8
	MsgID_BroadCast    = 9
	MsgID_PlayIPCAudio = 10
	MsgID_IPCControl   = 11

	MsgID_SetVolume = 13
)

const (
	MsgID_Device_Play_Audio  = 20
	MsgID_Device_Push_Audio  = 21
	MsgID_Device_Set_Volume  = 22
	MsgID_Device_Close_Audio = 23
)

var CC = map[string]int{
	StatusSucc:      http.StatusOK,
	StatusDBERR:     http.StatusServiceUnavailable,
	StatusParamsERR: http.StatusBadRequest,
	StatusAuthERR:   http.StatusUnauthorized,
	StatusSysERR:    http.StatusInternalServerError,
}

type Response struct {
	Data  any    `json:"data"`
	MsgID string `json:"msgid"`
	Code  string `json:"code"`
}

func JsonResponse(c *gin.Context, code string, data any) {
	logger.Logger.Debug("JsonResponse", zap.Any("code", code), zap.Any("data", data))
	switch d := data.(type) {
	case error:
		data = d.Error()
	}
	c.JSON(CC[code], Response{MsgID: c.GetString("msgid"), Code: code, Data: data})
}

func ZlmWebHookResponse(c *gin.Context, code int, msg string) {
	logger.Logger.Debug("ZlmWebHookResponse", zap.Any("code", code), zap.Any("msg", msg))
	statusCode := http.StatusOK
	if code == 401 {
		statusCode = http.StatusUnauthorized
	}
	c.JSON(statusCode, map[string]any{
		"code": code,
		"msg":  msg,
	})
}

const (
	DeviceStatusON  = "ON"
	DeviceStatusOFF = "OFF"
	defaultLimit    = 20
	defaultSort     = "-addtime"
)

func GetLimit(c *gin.Context) int {
	value := c.Query("limit")
	if value == "" {
		return defaultLimit
	}
	if d, e := strconv.Atoi(value); e == nil {
		return d
	} else {
		return defaultLimit
	}
}
func GetSort(c *gin.Context) string {
	value := c.Query("sort")
	if value == "" {
		return defaultSort
	}
	return value
}
func GetSkip(c *gin.Context) int {
	value := c.Query("skip")
	if value == "" {
		return 0
	}
	if d, e := strconv.Atoi(value); e == nil {
		return d
	} else {
		return 0
	}
}
