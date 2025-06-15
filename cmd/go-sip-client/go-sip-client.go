package main

import (
	capi "go-sip/api/c"
	. "go-sip/common"
	grpc_client "go-sip/grpc_api/c"
	"go-sip/logger"
	. "go-sip/logger"
	"go-sip/m"
	sipapi "go-sip/sip"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	m.LoadClientConfig()
	logger.InitLogger()

	device_id, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	r := gin.Default()

	r.POST(ZLMWebHookClientURL, capi.ZLMWebHook)

	client := grpc_client.NewSipClient(device_id)
	sipapi.Start()

	// 启动 API 服务，放到 goroutine 中
	go func() {
		for {
			err := r.Run(m.CMConfig.API)
			if err != nil {
				Logger.Error("api服务启动失败", zap.Error(err))
				time.Sleep(5 * time.Second)
				continue
			}
			Logger.Warn("api服务退出，尝试重新启动...")
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		tcp_addr := ""
		if tcp_addr == "" {
			// 默认获取配置中的tcp地址
			tcp_addr = m.CMConfig.TCP
		}
		if err := client.Connect(tcp_addr); err != nil {
			Logger.Error("连接失败", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		Logger.Info("连接成功")
		client.Run()
		Logger.Info("连接断开，尝试重新连接...")
	}
}
