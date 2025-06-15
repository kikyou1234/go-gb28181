package api

import (
	"encoding/json"
	"go-sip/db/redis"
	"go-sip/grpc_api"
	grpc_server "go-sip/grpc_api/s"
	"go-sip/m"
	pb "go-sip/signaling"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeviceControl(c *gin.Context) {
	ipc_id := c.Query("ipc_id")
	lr := c.Query("leftRight")
	ud := c.Query("upDown")
	io := c.Query("inOut")
	ms := c.Query("moveSpeed")
	zs := c.Query("zoomSpeed")

	_lr, err := strconv.Atoi(lr)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "leftRight 参数错误")
		return
	}
	_ud, err := strconv.Atoi(ud)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "upDown 参数错误")
		return
	}
	_io, err := strconv.Atoi(io)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "inOut 参数错误")
		return
	}
	_ms, err := strconv.Atoi(ms)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "moveSpeed 参数错误")
		return
	}
	_zs, err := strconv.Atoi(zs)
	if err != nil {
		m.JsonResponse(c, m.StatusParamsERR, "zoomSpeed 参数错误")
		return
	}

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_IPC_Control_Req{
		DeviceID:  ipc_id,
		LeftRight: _lr,
		UpDown:    _ud,
		InOut:     _io,
		MoveSpeed: _ms,
		ZoomSpeed: _zs,
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
		MsgID:   m.MsgID_IPCControl,
		Method:  m.DeviceControl,
		Payload: d,
	})
	if err != nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
		return
	}

	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}
