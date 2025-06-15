package grpc_client

import (
	"context"
	"fmt"
	db "go-sip/db/sqlite"
	"go-sip/grpc_api"
	. "go-sip/logger"
	"go-sip/m"
	pb "go-sip/signaling"
	sipapi "go-sip/sip"
	"go-sip/utils"

	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var SipCli *SipClient

var Stream_register = make(map[string]chan struct{})

type SipClient struct {
	clientID  string
	conn      *grpc.ClientConn
	stream    pb.SipService_StreamChannelClient
	client    pb.SipServiceClient
	AudioDone context.Context
}

func NewSipClient(clientID string) *SipClient {

	cl := &SipClient{
		clientID: clientID,
	}

	sipapi.NotifyFunc = cl.IpcEventReq
	sipapi.InviteFunc = cl.IpcInviteReq

	return cl
}

func GetSipClient() *SipClient {
	return SipCli
}

func (c *SipClient) Connect(addr string) error {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second, // PING发送间隔从默认15秒改为30秒
			Timeout:             20 * time.Second, // 等待PING响应的超时时间
			PermitWithoutStream: true,             // 允许无活跃流时发送PING
		}),
	)
	if err != nil {
		return err
	}

	client := pb.NewSipServiceClient(conn)

	stream, err := client.StreamChannel(context.Background())
	if err != nil {
		conn.Close()
		return err
	}

	// 发送注册信息
	if err := stream.Send(&pb.ClientMessage{
		Content: &pb.ClientMessage_Register{
			Register: &pb.ClientRegister{
				ClientId: c.clientID,
				Version:  "1.0.0",
			},
		},
	}); err != nil {
		conn.Close()
		return err
	}

	c.conn = conn
	c.stream = stream
	c.client = client
	return nil
}

func (c *SipClient) IpcEventReq(ipc_id string, msg_type string) (*pb.IpcEventAck, error) {

	if c.client != nil {
		return c.client.IpcEventReq(context.Background(), &pb.IpcEventRequest{
			ClientId: c.clientID,
			IpcId:    ipc_id,
			Event:    msg_type,
		})
	}

	return nil, fmt.Errorf("client以断开")

}

func (c *SipClient) IpcInviteReq(ipc_id string) (*pb.IpcInviteAck, error) {

	if c.client != nil {
		return c.client.IpcInviteReq(context.Background(), &pb.IpcInviteRequest{
			ClientId: c.clientID,
			IpcId:    ipc_id,
		})
	}

	return nil, fmt.Errorf("client以断开")

}

func (c *SipClient) Run() {
	defer c.conn.Close()

	// 用一个 channel 缓冲发送结果，避免多个 goroutine 并发写 stream
	resultChan := make(chan *pb.ClientMessage, 100)

	// 独立 goroutine 负责串行发送
	go func() {
		for msg := range resultChan {
			if err := c.stream.Send(msg); err != nil {
				Logger.Error("发送结果失败", zap.Error(err))
				return
			}
		}
	}()

	// 命令处理循环
	for {
		cmd, err := c.stream.Recv()
		if err != nil {
			Logger.Error("连接断开: ", zap.Error(err))
			close(resultChan) // 关闭发送协程
			return
		}

		go func(cmd *pb.ServerCommand) {
			// 处理服务端命令
			Logger.Info("收到服务端命令", zap.Any("method", cmd.Method), zap.Any("MsgID", cmd.MsgID))
			result := c.executeCommand(cmd)

			// 返回执行结果
			resultChan <- &pb.ClientMessage{
				Content: &pb.ClientMessage_Result{
					Result: &pb.CommandResult{
						MsgID:   cmd.MsgID,
						Success: result.Success,
						Payload: result.Msg,
					},
				},
			}

		}(cmd)
	}
}

type CommandResult struct {
	Success bool
	Msg     []byte
}

func (c *SipClient) executeCommand(cmd *pb.ServerCommand) CommandResult {

	rsp := CommandResult{}
	rsp.Success = true
	rsp.Msg = []byte("执行成功")
	switch cmd.Method {
	case m.Ping:

		return rsp

	case m.Play:
		d := &grpc_api.Sip_Play_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("实时流点播执行失败: %v", err))
			return rsp
		}

		// 向摄像头发送信令请求推实时流到zlm
		pm := &sipapi.Streams{ChannelID: d.ChannelID, StreamID: d.ChannelID,
			ZlmIP: d.ZLMIP, ZlmPort: d.ZLMPort, T: 0, Resolution: d.Resolution,
			Mode: 0, Ttag: db.M{}, Ftag: db.M{}, OnlyAudio: false}
		_, err = sipapi.SipPlay(pm)
		if err != nil {
			Logger.Error("向摄像头发送信令请求实时流推流到zlm失败", zap.Any("deviceId", d.ChannelID), zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte("向摄像头发送信令请求实时流推流到zlm失败")
			return rsp
		}

	case m.PlayBack:

		d := &grpc_api.Sip_Play_Back_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}
		s := time.Unix(d.StartTime, 0)
		e := time.Unix(d.EndTime, 0)
		pm := &sipapi.Streams{ChannelID: d.ChannelID, StreamID: fmt.Sprintf("%s_%d_%d", d.ChannelID, d.StartTime, d.EndTime), ZlmIP: d.ZLMIP, ZlmPort: d.ZLMPort, S: s, E: e, T: 1, Resolution: d.Resolution, Mode: d.Mode, Ttag: db.M{}, Ftag: db.M{}, OnlyAudio: false}
		_, err = sipapi.SipPlay(pm)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

	case m.StopPlay:
		d := &grpc_api.Sip_Stop_Play_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipStopPlay(d.StreamID)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}
	case m.ResumePlay:
		d := &grpc_api.Sip_Resume_Play_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipResumePlay(d.StreamID)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}
	case m.PausePlay:
		d := &grpc_api.Sip_Pause_Play_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipPausePlay(d.StreamID)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

	case m.SeekPlay:
		d := &grpc_api.Sip_Seek_Play_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipSeekPlay(d.StreamID, d.SubTime)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}
	case m.SpeedPlay:
		d := &grpc_api.Sip_Speed_Play_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipSpeedPlay(d.StreamID, d.Speed)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

	case m.RecordList:

		{

			d := &grpc_api.Sip_Play_Back_Recocd_List_Req{}
			err := utils.JSONDecode(cmd.Payload, d)
			if err != nil {
				Logger.Error("Unmarshal failed ", zap.Error(err))
				rsp.Success = false
				rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
				return rsp
			}

			channel := &sipapi.Channels{ChannelID: d.ChannelID}

			if err := db.Get(db.DBClient, channel); err != nil {
				rsp.Success = false
				rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
				return rsp
			}
			res, err := sipapi.SipRecordList(channel, d.StartTime, d.EndTime)
			if err != nil {
				rsp.Success = false
				rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
				return rsp
			}

			record := utils.JSONEncode(res)

			rsp.Msg = record
			return rsp
		}

	case m.Broadcast:

		d := &grpc_api.Sip_Ipc_BroadCast_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipIpcBroadCast(d.ChannelID)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}
	case m.PlayIPCAudio:
		d := &grpc_api.Sip_Play_IPC_Audio_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.SipPlayAudio(d.ChannelID, d.ZLMPort, d.ZLMIP)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

	case m.DeviceControl:
		d := &grpc_api.Sip_IPC_Control_Req{}
		err := utils.JSONDecode(cmd.Payload, d)
		if err != nil {
			Logger.Error("Unmarshal failed ", zap.Error(err))
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

		err = sipapi.DeviceControl(d.DeviceID, d.LeftRight, d.UpDown, d.InOut, d.MoveSpeed, d.ZoomSpeed)
		if err != nil {
			rsp.Success = false
			rsp.Msg = []byte(fmt.Sprintf("执行失败: %v", err))
			return rsp
		}

	}

	return rsp
}
