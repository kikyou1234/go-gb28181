package model

type ZlmInfo struct {
	ZlmIp     string `json:"zlmIp" validate:"required"`     // ZLM IP
	ZLMID     string `json:"zlmDomain" validate:"required"` // ZLM ID
	ZlmSecret string `json:"zlmSecret" validate:"required"` // ZLM 密钥
	ZlmPort   string `json:"zlmPort" validate:"required"`   // ZLM 端口
}
