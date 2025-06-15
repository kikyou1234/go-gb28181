package model

type ZLMStreamChangedData struct {
	Regist bool   `json:"regist"`
	Params string `json:"params"`
	APP    string `json:"app"`
	Stream string `json:"stream"`
	Schema string `json:"schema"`
}
