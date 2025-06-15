package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-sip/db/redis"
	"go-sip/grpc_api"
	grpc_server "go-sip/grpc_api/s"
	. "go-sip/logger"
	"go-sip/m"
	"go-sip/model"
	pb "go-sip/signaling"
	"go-sip/utils"
	"go-sip/zlm_api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var audio_pull_map = make(map[string]string) // audio 板端推流是否成功注册
var ZLM_Node = make(map[string]model.ZlmInfo)
var stream_hd = make(map[string]int)
var StreamWait = make(map[string]chan struct{})

func ZLMWebHook(c *gin.Context) {
	method := c.Param("method")
	Logger.Info("server sip ZLMWebHook", zap.String("method", method))
	switch method {
	case "on_server_started":
		ZlmServerStart(c)
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
	case "on_stream_changed":
		// 流注册和注销通知
		zlmStreamChanged(c)
	default:
		m.ZlmWebHookResponse(c, 0, "success")
	}

}

func ZlmServerStart(c *gin.Context) {
	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}

	var req model.ZlmServerStartDate
	if err := utils.JSONDecode(data, &req); err != nil {
		fmt.Println("JSON解析失败：", err)
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	zlmInfo := model.ZlmInfo{
		ZLMID:     req.MediaServerId,
		ZlmIp:     req.ExternIP,
		ZlmSecret: req.APISecret,
		ZlmPort:   req.HTTPPort,
	}

	zlmInfoJsonBytes, err := json.Marshal(zlmInfo)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "json marshal error")
		return
	}
	err = redis.HSet(c.Copy(), redis.ZLM_Node, req.MediaServerId, string(zlmInfoJsonBytes))
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "redis error")
		return
	}

	m.ZlmWebHookResponse(c, 0, "zlm server start success")
}

// zlm点播业务
func ZlmStreamOnPlay(c *gin.Context) {

	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	var req = model.ZlmStreamOnPlayData{}
	if err := utils.JSONDecode(data, &req); err != nil {
		m.ZlmWebHookResponse(c, -1, "参数格式错误，json反序列化失败")
		return
	}

	// 获取参数列表
	params := make(map[string]string)
	paramsArray := strings.Split(req.Params, "&")

	for _, param := range paramsArray {
		tokens := strings.Split(param, "=")
		if len(tokens) != 2 {
			m.ZlmWebHookResponse(c, -1, "参数格式错误，json反序列化失败")
			return
		}
		params[tokens[0]] = tokens[1]
	}
	redisZlmInfo, err := redis.HGet(c.Copy(), redis.ZLM_Node, req.MediaServerID)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "redis error")
		return
	}

	// 反序列化 JSON 字符串
	var zlmInfo model.ZlmInfo
	err = json.Unmarshal([]byte(redisZlmInfo), &zlmInfo)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "参数格式错误，json反序列化失败")
		return
	}

	stream_id_split := strings.Split(req.Stream, "_")
	// 实时流
	if len(stream_id_split) == 1 {
		ipc_id := stream_id_split[0]
		if req.App == "rtp" {
			hd, err := strconv.Atoi(params["hd"])
			if err != nil {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，hd 参数错误")
				return
			}

			mode, err := strconv.Atoi(params["mode"])
			if err != nil || mode < 0 || mode > 1 {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，mode 参数错误")
				return
			}
			device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
			if err != nil || device_id == "" {
				m.ZlmWebHookResponse(c, -1, "ipc_id未注册，请检查摄像头是否正常")
				return
			}

			if old_hd, ok := stream_hd[req.Stream]; ok {
				if old_hd != hd { // 切码流先停止点播
					sip_server := grpc_server.GetSipServer()
					sip_req := &grpc_api.Sip_Stop_Play_Req{
						StreamID: req.Stream,
					}
					d, err := json.Marshal(sip_req)
					if err != nil {
						m.ZlmWebHookResponse(c, -1, "参数格式错误，json序列化失败")
					}
					device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, req.Stream)
					if err != nil {
						m.ZlmWebHookResponse(c, -1, "ipc_id未注册，请检查摄像头是否正常")
					}
					_, err = sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
						MsgID:   m.MsgID_StopPlay,
						Method:  m.StopPlay,
						Payload: d,
					})
					if err != nil {
						m.ZlmWebHookResponse(c, -1, "终端请求错误，请检查是否掉线")
						return
					}
					zlm_api.ZlmCloseRtpServer("http://"+zlmInfo.ZlmIp+":"+zlmInfo.ZlmPort, zlmInfo.ZlmSecret, req.Stream)
				}
			}
		}
	} else if len(stream_id_split) == 3 { // 回放流
		Logger.Info("ZlmStreamOnPlay 回放流", zap.Any("req.Stream", req.Stream))
	} else {
		m.ZlmWebHookResponse(c, 403, "stream id 格式错误,实时流为 channelid  回放流为 channelid_starttime_endtime")
		return
	}
	//视频播放触发鉴权
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
		Logger.Info("流注册 ", zap.Any("req", req))

		if req.Schema == "rtsp" { // 如果是国标流注册
			if st, ok := StreamWait[req.Stream]; ok {
				st <- struct{}{}
			}
		}

	} else {

		Logger.Info("流注销 :", zap.Any("req", req))

		if req.APP == "audio" && req.Schema == "rtsp" {
			delete(audio_pull_map, req.Stream)
			delete(ZLM_Node, req.Stream)
			delete(stream_hd, req.Stream)
		}
	}
	m.ZlmWebHookResponse(c, 0, "success")
}

func zlmStreamPublish(c *gin.Context) {
	body := c.Request.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}
	var req = model.ZlmStreamPublishData{}
	if err := utils.JSONDecode(data, &req); err != nil {
		m.ZlmWebHookResponse(c, -1, "body error")
		return
	}

	// 获取参数列表
	params := make(map[string]string)
	paramsArray := strings.Split(req.Params, "&")

	if req.Params != "" {
		for _, param := range paramsArray {

			tokens := strings.Split(param, "=")

			if len(tokens) != 2 {
				m.ZlmWebHookResponse(c, 401, "Unauthorized")
				return
			}
			params[tokens[0]] = tokens[1]
		}
	}

	if req.App == "rtp" { // 国标推流 不用鉴权

	} else if req.App == "broadcast" { // 摄像头广播推流

		sip_req := &grpc_api.Sip_Ipc_BroadCast_Req{
			ChannelID: req.Stream,
		}
		d, err := json.Marshal(sip_req)
		if err != nil {
			m.ZlmWebHookResponse(c, -1, "参数格式错误，json序列化失败")
		}
		device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, req.Stream)
		if err != nil || device_id == "" {
			m.ZlmWebHookResponse(c, -1, "ipc_id未注册，请检查摄像头是否正常")
			return
		}

		sip_server := grpc_server.GetSipServer()

		sip_server.StreamMap[req.Stream] = req.MediaServerID
		go func() {

			_, err = sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
				MsgID:   m.MsgID_BroadCast,
				Method:  m.Broadcast,
				Payload: d,
			})
			if err != nil {
				return
			}
		}()

	} else {

		if sign, ok := params["sign"]; ok {
			if sign != utils.GetMD5(m.SMConfig.Sign) {
				m.ZlmWebHookResponse(c, 401, "Unauthorized")
				return
			}
		} else {
			m.ZlmWebHookResponse(c, 401, "Unauthorized")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"code":         0,
		"enable_audio": true,
		"enable_MP4":   false,
		"msg":          "success",
	})
}

// hd 0 标清 1 高清
// mode 0 UDP 1 Tcp被动
// http://127.0.0.1/index/api/webrtc?app=rtp&stream=37070000081118000001&type=play&hd=0&mode=0
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

	// 获取参数列表
	paramsMap := make(map[string]string)
	paramsArray := strings.Split(req.Params, "&")

	for _, param := range paramsArray {

		tokens := strings.Split(param, "=")
		if len(tokens) != 2 {
			m.ZlmWebHookResponse(c, -1, "传参格式错误")
			return
		}
		paramsMap[tokens[0]] = tokens[1]
	}
	redisZlmInfo, err := redis.HGet(c.Copy(), redis.ZLM_Node, req.MediaServerID)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "查询 ZLM Node 错误")
		return
	}

	// 反序列化 JSON 字符串
	var zlmInfo model.ZlmInfo
	err = json.Unmarshal([]byte(redisZlmInfo), &zlmInfo)
	if err != nil {
		m.ZlmWebHookResponse(c, -1, "参数格式错误，json序列化失败")
		return
	}
	sip_server := grpc_server.GetSipServer()
	stream_id_list := strings.Split(req.Stream, "_")

	if req.APP == "rtp" {
		device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, stream_id_list[0])
		if err != nil || device_id == "" {
			m.ZlmWebHookResponse(c, -1, "ipc_id未注册，请检查摄像头是否正常")
			return
		}

		mode, err := strconv.Atoi(paramsMap["mode"]) // 返回 (int, error)
		if err != nil || mode < 0 || mode > 1 {
			m.ZlmWebHookResponse(c, -1, "参数格式错误，mode 参数错误")
			return
		}

		resolution := 1
		if res, ok := paramsMap["hd"]; ok { // hd : 0 标清  1 高清
			if res == "1" {
				resolution = 1
			} else if res == "0" {
				resolution = 0
			} else {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，hd 参数错误")
				return
			}
		} else {
			m.ZlmWebHookResponse(c, -1, "参数格式错误，hd 参数错误")
			return
		}

		stream_hd[req.Stream] = resolution

		master_stream_id := stream_id_list[0]
		rtp_info := zlm_api.ZlmOpenRtpServer("http://"+zlmInfo.ZlmIp+":"+zlmInfo.ZlmPort, zlmInfo.ZlmSecret, req.Stream, mode)

		// rtp实时流
		if len(stream_id_list) == 1 {
			// 点播主ipc
			sip_req := &grpc_api.Sip_Play_Req{
				ChannelID:  master_stream_id,
				ZLMIP:      zlmInfo.ZlmIp,
				ZLMPort:    rtp_info.Port,
				Resolution: resolution,
				Mode:       mode,
			}
			d, err := json.Marshal(sip_req)
			if err != nil {
				Logger.Error("主ipc的json序列化失败")
			} else {
				_, err = sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
					MsgID:   m.MsgID_Play,
					Method:  m.Play,
					Payload: d,
				})
				if err != nil {
					Logger.Error("主ipc点播失败", zap.Any("ipcId", master_stream_id))
				}
			}

		} else if len(stream_id_list) == 3 { // 回放流
			device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, stream_id_list[0])
			if err != nil || device_id == "" {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，ipc_id未注册，请检查摄像头是否正常")
				return
			}

			startTime, err := strconv.ParseInt(stream_id_list[1], 10, 64)
			if err != nil {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，回放流开始时间转Int64失败")
				return
			}

			endTime, err := strconv.ParseInt(stream_id_list[2], 10, 64)
			if err != nil {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，回放流结束时间转Int64失败")
				return
			}

			// 点播回放流
			sip_req := &grpc_api.Sip_Play_Back_Req{
				ChannelID:  master_stream_id,
				ZLMIP:      zlmInfo.ZlmIp,
				ZLMPort:    rtp_info.Port,
				Resolution: resolution,
				StartTime:  startTime,
				EndTime:    endTime,
			}

			d, err := json.Marshal(sip_req)
			if err != nil {
				m.ZlmWebHookResponse(c, -1, "参数格式错误，json序列化失败")
				return
			}
			_, err = sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
				MsgID:   m.MsgID_PlayBack,
				Method:  m.PlayBack,
				Payload: d,
			})
			if err != nil {
				m.ZlmWebHookResponse(c, -1, "终端请求错误，请检查是否掉线")
				return
			}

		} else {
			m.ZlmWebHookResponse(c, -1, "stream id 格式错误,实时流为 channelid  回放流为 channelid_starttime_endtime")
			return
		}

	}
	StreamWait[req.Stream] = make(chan struct{})

	tick := time.NewTicker(5 * time.Second)
	select {
	case <-StreamWait[req.Stream]:
		close(StreamWait[req.Stream])
		delete(StreamWait, req.Stream)
		break
	case <-tick.C:
		Logger.Warn("等待流超时", zap.Any("stream", req.Stream))
		break
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

	if req.APP == "rtp" {

		s_size := strings.Split(req.Stream, "_")

		ipc_id := s_size[0]

		sip_server := grpc_server.GetSipServer()
		sip_req := &grpc_api.Sip_Stop_Play_Req{
			StreamID: req.Stream,
		}
		d, err := json.Marshal(sip_req)
		if err != nil {
			m.ZlmWebHookResponse(c, -1, "参数格式错误，json序列化失败")
		}
		device_id, err := redis.HGet(c.Copy(), redis.IPC_DEVICE, ipc_id)
		if err != nil {
			m.ZlmWebHookResponse(c, -1, "ipc_id未注册，请检查摄像头是否正常")
		}

		_, err = sip_server.ExecuteCommand(device_id, &pb.ServerCommand{
			MsgID:   m.MsgID_StopPlay,
			Method:  m.StopPlay,
			Payload: d,
		})
		if err != nil {
			m.ZlmWebHookResponse(c, -1, "终端请求错误，请检查是否掉线")
		}

	}
	c.JSON(http.StatusOK, map[string]any{
		"code":  0,
		"close": true,
	})
}
