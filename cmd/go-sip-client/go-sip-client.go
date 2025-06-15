package main

import (
	"encoding/json"
	"fmt"
	grpc_client "go-sip/grpc_api/c"
	"go-sip/logger"
	. "go-sip/logger"
	"go-sip/m"
	sipapi "go-sip/sip"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type SipServerResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func GetSipServer() (string, error) {
	url := fmt.Sprintf("http://%s/open/server/getone", m.CMConfig.Gateway) // 替换为你的实际接口地址

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 JSON
	var result SipServerResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Data, nil
}

func main() {

	m.LoadClientConfig()
	logger.InitLogger()

	device_id, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	client := grpc_client.NewSipClient(device_id)
	sipapi.Start()

	for {
		tcp_addr, err := GetSipServer()
		if err != nil {
			Logger.Error("连接sip-gateway失败", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		if err := client.Connect(tcp_addr); err != nil {
			Logger.Error("连接sip-server失败", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		Logger.Info("连接成功")
		client.Run()
		Logger.Info("连接断开，尝试重新连接...")
	}
}
