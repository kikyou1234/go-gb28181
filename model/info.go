package model

// SIP服务器信息结构体
type SipServerInfo struct {
	ServerID string `json:"server_id"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	SipIP    string `json:"sip_ip"`
	SipPort  string `json:"sip_port"`
}
