package common

// 定义Init中的url常量
const (
	// 设备列表接口
	DevicesListURL = "/devices"
	// 播放接口
	PlayURL     = "/ipc/play"
	PlaybackURL = "/ipc/playback"
	PauseURL    = "/ipc/pause"
	ResumeURL   = "/ipc/resume"
	SpeedURL    = "/ipc/speed"
	SeekURL     = "/ipc/seek"
	StopURL     = "/ipc/stop"

	DeviceControl = "/ipc/control"
	// 录像列表接口
	RecordsListURL = "/ipc/records"
	// ZLM Webhook接口
	ZLMWebHookBaseURL   = "/zlm/webhook"
	ZLMWebHookServerURL = ZLMWebHookBaseURL + "/:method"
	ZLMWebHookClientURL = "/client" + ZLMWebHookBaseURL + "/:method"
)
