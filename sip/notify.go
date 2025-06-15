package sipapi

import (
	. "go-sip/logger"
	"go-sip/signaling"

	"go.uber.org/zap"
)

const (
	// NotifyMethodUserActive 设备活跃状态通知
	NotifyMethodDevicesActive = "devices.active"
	// NotifyMethodUserRegister 设备注册通知
	NotifyMethodDevicesRegister = "devices.regiester"
	// NotifyMethodDeviceActive 通道活跃通知
	NotifyMethodChannelsActive = "channels.active"
	// INVITE_METHOD 通知方法
	NotifyMethodInvite = "invite"

	Local_ZLM_Host = "http://127.0.0.1:9092"
)

var NotifyFunc func(ipc_id string, msg_type string) (*signaling.IpcEventAck, error)
var InviteFunc func(ipc_id string) (*signaling.IpcInviteAck, error)

// Notify 消息通知结构
type Notify struct {
	Method   string `json:"method"`
	DeviceID string `json:"data"`
}

func notify(data *Notify) {

	defer func() {
		if r := recover(); r != nil {
			Logger.Error("notify panic recovered", zap.Any("error", r))
		}
	}()

	if data == nil {
		Logger.Error("notify received nil data")
		return
	}

	var ack *signaling.IpcEventAck
	var err error

	switch data.Method {
	case NotifyMethodDevicesActive: // 设备活跃通知
		ack, err = NotifyFunc(data.DeviceID, data.Method)
		if err != nil {
			Logger.Error("设备活跃通知服务端错误", zap.Any("data.DeviceID", data.DeviceID), zap.Error(err))
		}
		if ack == nil {
			Logger.Error("通知服务端返回ack为nil", zap.Any("data.DeviceID", data.DeviceID), zap.String("method", data.Method))
			return
		}
	case NotifyMethodDevicesRegister: // 设备注册通知
		ack, err = NotifyFunc(data.DeviceID, data.Method)
		if err != nil {
			Logger.Error("设备注册通知服务端错误", zap.Any("data.DeviceID", data.DeviceID), zap.Error(err))
		}
		if ack == nil {
			Logger.Error("通知服务端返回ack为nil", zap.Any("data.DeviceID", data.DeviceID), zap.String("method", data.Method))
			return
		}
	case NotifyMethodChannelsActive: // 通道活跃通知
		ack, err = NotifyFunc(data.DeviceID, data.Method)
		if err != nil {
			Logger.Error("通道活跃通知服务端错误", zap.Any("data.DeviceID", data.DeviceID), zap.Error(err))
		}
		if ack == nil {
			Logger.Error("通知服务端返回ack为nil", zap.Any("data.DeviceID", data.DeviceID), zap.String("method", data.Method))
			return
		}
		// 收到摄像头发出的通知，调用内部zlm接口，查询流id是否存在
		// var req = zlm_api.ZlmGetRtpInfoReq{}
		// req.App = "rtp"
		// req.StreamID = data.DeviceID
		// req.Vhost = "__defaultVhost__"
		// resp := zlm_api.ZlmGetRtpInfo(Local_ZLM_Host, m.CMConfig.ZlmSecret, req)
		// if resp.Code == 0 && !resp.Exist {
		// 	Logger.Info("本地ZLM不存在流，开始推送本地流", zap.Any("deviceID", data.DeviceID))
		// 	// 不存在，调用openRtpServer
		// 	rtp_info := zlm_api.ZlmOpenRtpServer(Local_ZLM_Host, m.CMConfig.ZlmSecret, data.DeviceID, 0)
		// 	if rtp_info.Code != 0 || rtp_info.Port == 0 {
		// 		Logger.Error("open rtp server fail", zap.Int("code", rtp_info.Code))
		// 		return
		// 	}
		// 	// 向摄像头发送信令请求推实时流到zlm
		// 	pm := &Streams{ChannelID: data.DeviceID, StreamID: data.DeviceID,
		// 		ZlmIP: m.CMConfig.ZlmInnerIp, ZlmPort: rtp_info.Port, T: 0, Resolution: 0,
		// 		Mode: 0, Ttag: db.M{}, Ftag: db.M{}, OnlyAudio: false}
		// 	_, err = SipPlay(pm)
		// 	if err != nil {
		// 		Logger.Error("向摄像头发送信令请求实时流推流到zlm失败", zap.Any("deviceId", data.DeviceID), zap.Error(err))
		// 		return
		// 	}
		// }
	case NotifyMethodInvite: // 邀请通知
		ack, err = NotifyFunc(data.DeviceID, data.Method)
		if err != nil {
			Logger.Error("邀请通知服务端错误", zap.Any("data.DeviceID", data.DeviceID), zap.Error(err))
		}
		if ack == nil {
			Logger.Error("通知服务端返回ack为nil", zap.Any("data.DeviceID", data.DeviceID), zap.String("method", data.Method))
			return
		}
	default:
		Logger.Error("notify config not found", zap.Any("data.DeviceID", data.DeviceID), zap.Any("method", data.Method))
	}

}

func notifyDevicesAcitve(device_id string) *Notify {
	return &Notify{
		Method:   NotifyMethodDevicesActive,
		DeviceID: device_id,
	}
}
func notifyDevicesRegister(device_id string) *Notify {
	return &Notify{
		Method:   NotifyMethodDevicesRegister,
		DeviceID: device_id,
	}
}

func notifyChannelsActive(channelid string) *Notify {
	return &Notify{
		Method:   NotifyMethodChannelsActive,
		DeviceID: channelid,
	}
}

func notifyInvite(channelid string) *Notify {
	return &Notify{
		Method:   NotifyMethodInvite,
		DeviceID: channelid,
	}
}
