package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"go-sip/api"
	"go-sip/api/middleware"
	"go-sip/db/redis"
	grpc_server "go-sip/grpc_api/s"
	"go-sip/logger"
	"go-sip/m"
	"go-sip/model"
	pb "go-sip/signaling"
	"go-sip/utils"

	. "go-sip/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func RegisterGateway() error {

	t := model.SipServerInfo{
		ServerID: m.SMConfig.SipID,
		IP:       m.SMConfig.Api.IP,
		Port:     m.SMConfig.Api.Port,
	}

	data := utils.JSONEncode(t)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/open/server/register", m.SMConfig.GatewayAPI), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("注册失败")
	}

	return nil

}

func main() {
	m.LoadServerConfig()
	logger.InitLogger()

	redis.InitRedis(m.SMConfig.DateBase.Host, m.SMConfig.DateBase.Password, m.SMConfig.DateBase.DB)
	r := gin.Default()
	r.Use(middleware.Recovery)

	api.ServerApiInit(r)

	err := RegisterGateway()
	if err != nil {
		panic(err)
	}
	lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%s", m.SMConfig.Sip.SipIP, m.SMConfig.Sip.SipPort))
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 15 * time.Minute, // 空闲连接最大保持时间
			Time:              30 * time.Second, // PING发送间隔（服务端->客户端）
			Timeout:           20 * time.Second, // PING响应超时
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second, // 允许客户端的最小PING间隔
			PermitWithoutStream: true,             // 允许无活跃流的PING
		}),
	)
	sip_server := grpc_server.GetSipServer()
	pb.RegisterSipServiceServer(grpcServer, sip_server)
	go grpcServer.Serve(lis)

	err = r.Run(fmt.Sprintf("%s:%s", m.SMConfig.Api.IP, m.SMConfig.Api.Port))
	if err != nil {
		Logger.Error("sip server启动失败", zap.Error(err))
		return
	}
}
