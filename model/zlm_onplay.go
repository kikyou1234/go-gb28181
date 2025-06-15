package model

type ZlmStreamOnPlayData struct {
	App           string `json:"app"`
	HookIndex     int    `json:"hook_index"`
	ID            string `json:"id"`
	IP            string `json:"ip"`
	MediaServerID string `json:"mediaServerId"`
	Params        string `json:"params"`
	Port          int    `json:"port"`
	Protocol      string `json:"protocol"`
	Schema        string `json:"schema"`
	Stream        string `json:"stream"`
	Vhost         string `json:"vhost"`
}
