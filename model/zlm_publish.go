package model

type ZlmStreamPublishData struct {
	MediaServerID string `json:"mediaServerId"` // 媒体服务器ID
	App           string `json:"app"`           // 应用名称
	ID            string `json:"id"`            // 流唯一标识符
	IP            string `json:"ip"`            // 客户端IP地址
	Params        string `json:"params"`        // 附加参数
	Port          int    `json:"port"`          // 客户端端口
	Schema        string `json:"schema"`        // 协议方案
	Protocol      string `json:"protocol"`      // 实际使用协议
	Stream        string `json:"stream"`        // 流名称
	VHost         string `json:"vhost"`         // 虚拟主机
}
