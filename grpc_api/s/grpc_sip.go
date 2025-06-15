package grpc_server

import (
	"context"
	"encoding/json"

	"go-sip/db/redis"

	. "go-sip/logger"
	"go-sip/m"
	"go-sip/model"
	pb "go-sip/signaling"
	sipapi "go-sip/sip"
	"go-sip/zlm_api"

	"sync"
	"time"

	"github.com/gogo/status"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

var SipSrv *SipServer

type SipServer struct {
	pb.UnimplementedSipServiceServer
	clients   sync.Map // 使用 sync.Map 管理客户端连接
	StreamMap map[string]string
}

func GetSipServer() *SipServer {

	if SipSrv == nil {
		SipSrv = &SipServer{
			StreamMap: make(map[string]string),
			clients:   sync.Map{},
		}
	}
	return SipSrv
}

// 双向流处理
func (s *SipServer) StreamChannel(stream pb.SipService_StreamChannelServer) error {
	// 接收初始注册信息
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}

	reg := firstMsg.GetRegister()
	if reg == nil {
		return status.Error(codes.InvalidArgument, "客户端注册失败")
	}

	// 记录客户端连接
	clientCtx := &ClientContext{
		ID:     reg.ClientId,
		Stream: stream,
	}

	redis.HSet(context.TODO(), redis.SIPSERVER_DEVICE, reg.ClientId, m.SMConfig.SipID)

	s.clients.Store(reg.ClientId, clientCtx)
	defer s.clients.Delete(reg.ClientId)
	defer redis.HDel(context.TODO(), redis.SIPSERVER_DEVICE, reg.ClientId)

	for {
		msg, err := stream.Recv()
		if err != nil {
			Logger.Error("客户端连接已断开", zap.Any("client id", reg.ClientId), zap.Error(err))
			return err
		}
		if msg != nil {
			// 如果是响应
			if res := msg.GetResult(); res != nil {
				if chVal, ok := clientCtx.ResponseChans.Load(res.MsgID); ok {
					ch := chVal.(chan *pb.CommandResult)
					ch <- res
					close(ch)
					clientCtx.ResponseChans.Delete(res.MsgID)
				}
				continue
			}
		}
	}

}

// 主动调用客户端方法
func (s *SipServer) ExecuteCommand(clientID string, cmd *pb.ServerCommand) (*pb.CommandResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, ok := s.clients.Load(clientID)
	if !ok {
		return nil, status.Error(codes.NotFound, "客户端未连接")
	}
	client := val.(*ClientContext)

	resultChan := make(chan *pb.CommandResult, 1)
	client.ResponseChans.Store(cmd.MsgID, resultChan)
	defer client.ResponseChans.Delete(cmd.MsgID)

	if err := client.Stream.Send(cmd); err != nil {
		return nil, err
	}

	select {
	case res := <-resultChan:
		return res, nil
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "等待响应超时")
	}
}

type ClientContext struct {
	ID            string
	Stream        pb.SipService_StreamChannelServer
	LastActive    time.Time
	ClientCtx     context.Context
	ResponseChans sync.Map // key: MsgID(string), value: chan *pb.CommandResult
}

func (s *SipServer) IpcEventReq(ctx context.Context, req *pb.IpcEventRequest) (*pb.IpcEventAck, error) {

	switch req.Event {
	case sipapi.NotifyMethodDevicesActive:

	case sipapi.NotifyMethodDevicesRegister:
		Logger.Info("添加新的摄像头 ", zap.Any("device_id", req.ClientId), zap.Any("ipc_id", req.IpcId))
	case sipapi.NotifyMethodChannelsActive:
		Logger.Info("收到通道活跃通知 ", zap.Any("device_id", req.ClientId), zap.Any("channel_id", req.IpcId))
		redis.HSet(ctx, redis.IPC_SIPSERVER, req.IpcId, m.SMConfig.SipID) // ipc_id -> sip-server 记录 这个摄像头在 哪个sip-server
		redis.HSet(ctx, redis.IPC_DEVICE, req.IpcId, req.ClientId)        //  ipc-id -> client_id 记录 这个摄像头在哪个 客户端

	}
	return &pb.IpcEventAck{Success: true, Msg: "success"}, nil
}

func (s *SipServer) IpcInviteReq(ctx context.Context, req *pb.IpcInviteRequest) (*pb.IpcInviteAck, error) {

	zlm_id := s.StreamMap[req.IpcId]
	redisZlmInfo, err := redis.HGet(ctx, redis.ZLM_Node, zlm_id)
	if err != nil {
		return nil, err
	}
	// 反序列化 JSON 字符串
	var zlmInfo model.ZlmInfo
	err = json.Unmarshal([]byte(redisZlmInfo), &zlmInfo)
	if err != nil {
		return nil, err
	}

	rtp_info := zlm_api.ZlmStartSendRtpPassive(zlmInfo.ZlmIp, zlmInfo.ZlmPort, zlmInfo.ZlmSecret, req.IpcId)

	return &pb.IpcInviteAck{Success: true, ZlmIp: zlmInfo.ZlmIp, ZlmPort: int64(rtp_info.LocalPort)}, nil
}
