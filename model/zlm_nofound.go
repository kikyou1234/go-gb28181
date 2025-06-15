package model

type ZLMStreamNotFoundData struct {
	APP           string `json:"app"`
	Params        string `json:"params"`
	Stream        string `json:"stream"`
	Schema        string `json:"schema"`
	ID            string `json:"id"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	MediaServerID string `json:"mediaServerId"`
}
