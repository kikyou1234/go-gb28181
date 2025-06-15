package sipapi

import (
	"time"

	db "go-sip/db/sqlite"
	. "go-sip/logger"
	"go-sip/utils"

	"go.uber.org/zap"
)

// MessageNotify 心跳包xml结构
type MessageNotify struct {
	CmdType  string `xml:"CmdType"`
	SN       int    `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
	Status   string `xml:"Status"`
	Info     string `xml:"Info"`
}

func sipMessageKeepalive(u Devices, body []byte) error {
	message := &MessageNotify{}
	if err := utils.XMLDecode(body, message); err != nil {
		Logger.Error("Message Unmarshal xml err:", zap.Error(err))
		return err
	}
	device, ok := _activeDevices.Get(u.DeviceID)
	if !ok {
		device = Devices{DeviceID: u.DeviceID}
		if err := db.Get(db.DBClient, &device); err != nil {
			Logger.Error("device not found:", zap.Any("device id", u.DeviceID), zap.Error(err))
		}
	}
	if message.Status == "OK" {
		device.ActiveAt = time.Now().Unix()
		_activeDevices.Store(u.DeviceID, u)
	} else {
		device.ActiveAt = -1
		_activeDevices.Delete(u.DeviceID)
	}
	go notify(notifyDevicesAcitve(u.DeviceID))
	_, err := db.UpdateAll(db.DBClient, new(Devices), map[string]interface{}{"deviceid=?": u.DeviceID}, Devices{
		Host:     u.Host,
		Port:     u.Port,
		Rport:    u.Rport,
		RAddr:    u.RAddr,
		Source:   u.Source,
		URIStr:   u.URIStr,
		ActiveAt: device.ActiveAt,
	})
	return err
}
