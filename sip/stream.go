package sipapi

import (
	"fmt"
	"sync"
	"time"

	db "go-sip/db/sqlite"

	sip "go-sip/sip/s"
)

// Streams Streams
type Streams struct {
	db.DBModel
	// 0  直播 1 历史
	T int `json:"t" gorm:"column:t"`
	// 设备ID
	DeviceID string `json:"deviceid" gorm:"column:deviceid"`
	// 通道ID
	ChannelID string `json:"channelid" gorm:"column:channelid"`
	//  pull 媒体服务器主动拉流，push 监控设备主动推流
	StreamType string `json:"streamtype" gorm:"column:streamtype"`
	// 0正常 1关闭 -1 尚未开始
	Status int `json:"status" gorm:"column:status"`
	// header from params
	Ftag db.M `gorm:"column:ftag" sql:"type:json" json:"-"`
	// header to params
	Ttag db.M `gorm:"column:ttag" sql:"type:json" json:"-"`
	// header callid
	CallID string `json:"callid" gorm:"column:callid"`
	// 是否停止
	Stop   bool   `json:"stop" gorm:"column:stop"`
	Msg    string `json:"msg" gorm:"column:msg"`
	CseqNo uint32 `json:"cseqno" gorm:"column:cseqno"`
	// 视频流ID gb28181的ssrc
	StreamID string `json:"streamid"  gorm:"column:streamid"`
	// zlm是否收到流
	Stream bool `json:"stream" gorm:"column:stream"`

	ZlmIP string `json:"zlm_ip" gorm:"column:zlm_ip"`

	ZlmPort int `json:"zlm_port" gorm:"column:zlm_port"`

	// 0 标清，1 高清
	Resolution int `json:"-" gorm:"-"`

	OnlyAudio bool `json:"only_audio" gorm:"column:only_audio"` // 是否只拉音频流

	// 0 udp 模式发流 1 tcp 被动发流
	Mode int `json:"-" gorm:"-"`
	// ---
	S, E time.Time     `json:"-" gorm:"-"`
	ssrc string        // 国标ssrc 10进制字符串
	Ext  int64         `json:"-" gorm:"-"` // 流等待过期时间
	Resp *sip.Response `json:"-" gorm:"-"`
}

// 当前系统中存在的流列表
type streamsList struct {
	// key=ssrc value=PlayParams  播放对应的PlayParams 用来发送bye获取tag，callid等数据
	Response *sync.Map
	// key=channelid value={Play}  当前设备直播信息，防止重复直播
	Succ *sync.Map
	ssrc int
}

var StreamList streamsList

func getSSRC(t int) string {
	r := false
	for {
		StreamList.ssrc++
		// ssrc最大为四位数，超过时从1开始重新计算
		if StreamList.ssrc > 9000 && !r {
			StreamList.ssrc = 0
			r = true
		}
		key := fmt.Sprintf("%d%s%04d", t, _sysinfo.Region[3:8], StreamList.ssrc)
		stream := Streams{StreamID: ssrc2stream(key), Stop: false}
		if err := db.Get(db.DBClient, &stream); db.RecordNotFound(err) || stream.CreatedAt == 0 {
			return key
		}
	}
}
