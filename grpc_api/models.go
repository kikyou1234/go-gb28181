package grpc_api

type Sip_Play_Req struct {
	ChannelID  string
	ZLMIP      string
	ZLMPort    int
	Resolution int
	Mode       int
}

type Sip_Play_Back_Req struct {
	ChannelID  string
	ZLMIP      string // zlm ip
	ZLMPort    int    // zlm port
	Resolution int    // 0 标清 1 高清
	Mode       int    // 0 udp 1 tcp
	StartTime  int64
	EndTime    int64
}

type Sip_Stop_Play_Req struct {
	StreamID string
}

type Sip_Play_Back_Recocd_List_Req struct {
	ChannelID string
	StartTime int64
	EndTime   int64
}

type Sip_Pause_Play_Req struct {
	StreamID string
}

type Sip_Resume_Play_Req struct {
	StreamID string
}

type Sip_Seek_Play_Req struct {
	StreamID string
	SubTime  int64
}

type Sip_Speed_Play_Req struct {
	StreamID string
	Speed    float64
}

type Sip_Ipc_BroadCast_Req struct {
	ChannelID string
}

type Sip_Play_IPC_Audio_Req struct {
	ChannelID string
	ZLMIP     string
	ZLMPort   int
}

type Sip_IPC_Control_Req struct {
	DeviceID  string
	LeftRight int // 镜头左移右移 0:停止 1:左移 2:右移
	UpDown    int // 镜头上移下移 0:停止 1:上移 2:下移
	InOut     int // 镜头放大缩小 0:停止 1:缩小 2:放大
	MoveSpeed int // 镜头移动速度
	ZoomSpeed int // 镜头缩放速度
}
