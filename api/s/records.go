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

// @Summary     回放文件时间列表
// @Description 用来获取通道设备存储的可回放时间段列表，注意控制时间跨度，跨度越大，数据量越多，返回越慢，甚至会超时（最多10s）。
// @Tags        records
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Param       id    path     string true "通道id"
// @Param       start query    int    true "开始时间，时间戳"
// @Param       end   query    int    true "结束时间，时间戳"
// @Success     0     {object} sipapi.Records
// @Failure     1000  {object} string
// @Failure     1001  {object} string
// @Failure     1002  {object} string
// @Failure     1003  {object} string
// @Router      /channels/{id}/records [get]
func RecordsList(c *gin.Context) {
	ipc_id := c.Query("ipc_id")
	start := c.Query("start")
	end := c.Query("end")

	if start == "" {
		m.JsonResponse(c, m.StatusParamsERR, "开始时间错误")
		return
	}
	startStamp, err := strconv.ParseInt(start, 10, 64)
	if err != nil || startStamp <= 0 {
		m.JsonResponse(c, m.StatusParamsERR, "开始时间错误")
		return
	}
	if end == "" {
		m.JsonResponse(c, m.StatusParamsERR, "结束时间错误")
		return
	}
	endStamp, err := strconv.ParseInt(end, 10, 64)
	if err != nil || endStamp <= 0 || endStamp <= startStamp {
		m.JsonResponse(c, m.StatusParamsERR, "结束时间错误")
		return
	}

	sip_server := grpc_server.GetSipServer()
	data := &grpc_api.Sip_Play_Back_Recocd_List_Req{
		ChannelID: ipc_id,
		StartTime: startStamp,
		EndTime:   endStamp,
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
		MsgID:   m.MsgID_RecordList,
		Method:  m.RecordList,
		Payload: d,
	})
	if err != nil {
		m.JsonResponse(c, m.StatusSysERR, "中控请求错误，请检查是否掉线")
		return
	}

	m.JsonResponse(c, m.StatusSucc, string(result.Payload))
}
