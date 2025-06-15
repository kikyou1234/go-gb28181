package api

import (
	"encoding/json"
	"fmt"
	"go-sip/db/redis"
	"go-sip/grpc_api"
	grpc_server "go-sip/grpc_api/s"
	"go-sip/m"
	pb "go-sip/signaling"
	"go-sip/zlm_api"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary     监控播放（直播/回放）
// @Description 直播一个通道最多存在一个流，回放每请求一次生成一个流
// @Tags        streams
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       ipc_id     path     string true  "通道id"
// @Param       audio_type formData int    false "1高清，0标清，默认0"
// @Success     0      {object} sipapi.Streams
// @Failure     1000 {object} string
// @Failure     1001 {object} string
// @Failure     1002 {object} string
// @Failure     1003 {object} string
// @Router      /channels/{id}/streams [post]
func Play(c *gin.Context) {
	ipc_id := c.Query("ipc_id")

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Play_Req{
		ChannelID: ipc_id,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
	}
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}

	result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
		MsgID:   m.MsgID_Play,
		Method:  m.Play,
		Payload: d,
	})
	if err != nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
	}

	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}

// @Summary     监控回放
// @Description 直播一个通道最多存在一个流，回放每请求一次生成一个流
// @Tags        streams
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       id     path     string true  "通道id"
// @Param       start  formData int    false "回放开始时间，时间戳，replay=1时必传"
// @Param       end    formData int    false "回放结束时间，时间戳，replay=1时必传"
// @Success     0      {object} sipapi.Streams
// @Failure     1000 {object} string
// @Failure     1001 {object} string
// @Failure     1002 {object} string
// @Failure     1003 {object} string
// @Router      /channels/{id}/streams [post]
func Playback(c *gin.Context) {
	ipc_id := c.Param("ipc_id")

	s, _ := strconv.ParseInt(c.PostForm("start"), 10, 64)
	if s == 0 {
		m.JsonResponse(c, m.StatusParamsERR, "开始时间错误")
		return
	}
	e, _ := strconv.ParseInt(c.PostForm("end"), 10, 64)
	if e == 0 {
		m.JsonResponse(c, m.StatusParamsERR, "结束时间错误")
		return
	}
	if s >= e {
		m.JsonResponse(c, m.StatusParamsERR, "开始时间>=结束时间")
		return
	}

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Play_Back_Req{
		ChannelID: ipc_id,
		StartTime: s,
		EndTime:   e,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
	}
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}
	result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
		MsgID:   m.MsgID_PlayBack,
		Method:  m.PlayBack,
		Payload: d,
	})
	if err != nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
	}

	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}

// @Summary     停止播放（直播/回放）
// @Description 无人观看5分钟自动关闭，直播流无需调用此接口。
// @Tags        streams
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       ipc_id   path     string true "流id,播放接口返回的streamid"
// @Success     0    {object} string
// @Failure     1000   {object} string
// @Failure     1001   {object} string
// @Failure     1002   {object} string
// @Failure     1003   {object} string
// @Router      /streams/{id} [delete]
func Stop(c *gin.Context) {
	stream_id := c.Query("stream")

	s_size := strings.Split(stream_id, "_")
	if len(s_size) != 3 {
		m.JsonResponse(c, m.StatusParamsERR, "回放流ID错误")
		return
	}

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Stop_Play_Req{
		StreamID: stream_id,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
	}

	ipc_id := s_size[0]
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}

	result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
		MsgID:   m.MsgID_StopPlay,
		Method:  m.StopPlay,
		Payload: d,
	})
	if err != nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
		return
	}
	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}

// @Summary     暂停播放（直播/回放）
func Pause(c *gin.Context) {
	stream_id := c.Query("stream")
	zlm_ip := c.Query("zlmIp")
	zlm_secret := c.Query("zlmSecret")

	s_size := strings.Split(stream_id, "_")

	if len(s_size) != 3 {
		m.JsonResponse(c, m.StatusParamsERR, "回放流ID错误")
		return
	}

	ipc_id := s_size[0]

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Pause_Play_Req{
		StreamID: stream_id,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
		return
	}
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}

	rsp := zlm_api.ZlmPauseRtpCheck(fmt.Sprintf("http://%s:9092", zlm_ip), zlm_secret, stream_id)
	if rsp.Code == 0 {
		result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
			MsgID:   m.MsgID_PausePlay,
			Method:  m.PausePlay,
			Payload: d,
		})
		if err != nil {
			m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
			return
		}
		m.JsonResponse(c, m.StatusSucc, string(result.Payload))
		return
	} else {
		m.JsonResponse(c, m.StatusSysERR, "暂停失败")
		return
	}

}

func Resume(c *gin.Context) {
	stream_id := c.Query("stream")

	s_size := strings.Split(stream_id, "_")
	if len(s_size) != 3 {
		m.JsonResponse(c, m.StatusParamsERR, "回放流ID错误")
		return
	}

	ipc_id := s_size[0]

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Resume_Play_Req{
		StreamID: stream_id,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
		return
	}
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}
	rsp := zlm_api.ZlmResumeRtpCheck("http://127.0.0.1:9092", "ShAngHULasNduFxW681tYivExOLXaO3S", stream_id)

	if rsp.Code == 0 {
		result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
			MsgID:   m.MsgID_ResumePlay,
			Method:  m.ResumePlay,
			Payload: d,
		})
		if err != nil {
			m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
			return
		}
		m.JsonResponse(c, m.StatusSucc, string(result.Payload))
		return
	} else {
		m.JsonResponse(c, m.StatusSysERR, "暂停失败")
		return
	}
}

func Speed(c *gin.Context) {
	stream_id := c.Query("stream")
	speed := c.Query("speed")

	s, err := strconv.ParseFloat(speed, 64)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误")
		return
	}

	s_size := strings.Split(stream_id, "_")
	if len(s_size) != 3 {
		m.JsonResponse(c, m.StatusParamsERR, "回放流ID错误")
		return
	}

	ipc_id := s_size[0]

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Speed_Play_Req{
		StreamID: stream_id,
		Speed:    s,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
		return
	}
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}

	result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
		MsgID:   m.MsgID_SpeedPlay,
		Method:  m.SpeedPlay,
		Payload: d,
	})
	if err != nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
		return
	}
	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}

func Seek(c *gin.Context) {
	stream_id := c.Query("stream")
	seek := c.Query("seek")

	s, err := strconv.ParseInt(seek, 10, 64)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误,seek 必须为整型")
		return
	}

	s_size := strings.Split(stream_id, "_")
	if len(s_size) != 3 {
		m.JsonResponse(c, m.StatusParamsERR, "回放流ID错误")
		return
	}

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Seek_Play_Req{
		StreamID: stream_id,
		SubTime:  s,
	}
	d, err := json.Marshal(data)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "参数格式错误，json序列化失败")
		return
	}

	ipc_id := s_size[0]
	device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "ipc_id未注册，请检查摄像头是否正常")
		return
	}

	result, err := sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
		MsgID:   m.MsgID_SeekPlay,
		Method:  m.SeekPlay,
		Payload: d,
	})
	if err != nil || result == nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
		return
	}
	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}
