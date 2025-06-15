package model

type StreamMergeConfigDTO struct {
	GapV   int        `json:"gapv"`
	GapH   int        `json:"gaph"`
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Row    int        `json:"row"`
	Col    int        `json:"col"`
	ID     string     `json:"id"`
	URL    [][]string `json:"url"`
	Span   []int      `json:"span"`
}

type StreamMergeInfoDTO struct {
	DeviceId  string   `json:"deviceId" validate:"required"`        // 设备ID
	IpcIdList []string `json:"ipcIdList" validate:"required,min=1"` // 国标设备ID列表，不能为空
	StreamId  string   `json:"streamId" validate:"required"`        // 流 ID
	Type      int      `json:"type" validate:"required"`            // 合屏还是切屏 1 合屏 2 切屏
}
