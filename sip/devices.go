package sipapi

import (
	"encoding/xml"
	"fmt"
	"net"
	"time"

	db "go-sip/db/sqlite"
	"go-sip/m"
	sip "go-sip/sip/s"
	"go-sip/utils"

	. "go-sip/logger"

	"go.uber.org/zap"
)

var (
	// sip服务用户信息
	_serverDevices Devices
	srv            *sip.Server
)

// Devices NVR  设备信息
type Devices struct {
	db.DBModel
	// Name 设备名称
	Name string `json:"name" gorm:"column:name" `
	// DeviceID 设备id
	DeviceID string `json:"deviceid" gorm:"column:deviceid"`
	// Region 设备域
	Region string `json:"region" gorm:"column:region"`
	// Host Via 地址
	Host string `json:"host" gorm:"column:host"`
	// Port via 端口
	Port string `json:"port" gorm:"column:port"`
	// TransPort via transport
	TransPort string `json:"transport" gorm:"column:transport"`
	// Proto 协议
	Proto string `json:"proto" gorm:"column:proto"`
	// Rport via rport
	Rport string `json:"report" gorm:"column:report"`
	// RAddr via recevied
	RAddr string `json:"raddr"  gorm:"column:raddr"`
	// Manufacturer 制造厂商
	Manufacturer string `xml:"Manufacturer"  json:"manufacturer"  gorm:"column:manufacturer"`
	// 设备类型DVR，NVR
	DeviceType string `xml:"DeviceType"  json:"devicetype"  gorm:"column:devicetype"`
	// Firmware 固件版本
	Firmware string ` json:"firmware"  gorm:"column:firmware"`
	// Model 型号
	Model  string `json:"model"  gorm:"column:model"`
	URIStr string `json:"uri"  gorm:"column:uri"`
	// ActiveAt 最后心跳检测时间
	ActiveAt int64 `json:"active" gorm:"column:active"`
	// Regist 是否注册
	Regist bool `json:"regist"  gorm:"column:regist"`
	// PWD 密码
	PWD string `json:"pwd" gorm:"column:pwd"`
	// Source
	Source string `json:"source"  gorm:"column:source"`

	Sys m.SysInfo `json:"sysinfo" gorm:"-"`

	//----
	addr   *sip.Address `gorm:"-"`
	source net.Addr     `gorm:"-"`
}

// Channels 摄像头通道信息
type Channels struct {
	db.DBModel
	// ChannelID 通道编码
	ChannelID string `xml:"DeviceID" json:"channelid" gorm:"column:channelid"`
	// DeviceID 设备编号
	DeviceID string `xml:"-" json:"deviceid"  gorm:"column:deviceid"`
	// Memo 备注（用来标示通道信息）
	MeMo string `json:"memo"  gorm:"column:memo"`
	// Name 通道名称（设备端设置名称）
	Name         string `xml:"Name" json:"name"  gorm:"column:name"`
	Manufacturer string `xml:"Manufacturer" json:"manufacturer"  gorm:"column:manufacturer"`
	Model        string `xml:"Model" json:"model"  gorm:"column:model"`
	Owner        string `xml:"Owner"  json:"owner"  gorm:"column:owner"`
	CivilCode    string `xml:"CivilCode" json:"civilcode"  gorm:"column:civilcode"`
	// Address ip地址
	Address     string `xml:"Address"  json:"address"  gorm:"column:address"`
	Parental    int    `xml:"Parental"  json:"parental"  gorm:"column:parental"`
	SafetyWay   int    `xml:"SafetyWay"  json:"safetyway"  gorm:"column:safetyway"`
	RegisterWay int    `xml:"RegisterWay"  json:"registerway"  gorm:"column:registerway"`
	Secrecy     int    `xml:"Secrecy" json:"secrecy"  gorm:"column:secrecy"`
	// Status 状态  on 在线
	Status string `xml:"Status"  json:"status"  gorm:"column:status"`
	// Active 最后活跃时间
	Active int64  `json:"active"  gorm:"column:active"`
	URIStr string ` json:"uri"  gorm:"column:uri"`

	// 视频编码格式
	VF string ` json:"vf"  gorm:"column:vf"`
	// 视频高
	Height int `json:"height"  gorm:"column:height"`
	// 视频宽
	Width int `json:"width"  gorm:"column:width"`
	// 视频FPS
	FPS int `json:"fps"  gorm:"column:fps"`
	//  pull 媒体服务器主动拉流，push 监控设备主动推流
	StreamType string `json:"streamtype"  gorm:"column:streamtype"`
	// streamtype=pull时，拉流地址
	URL string `json:"url"  gorm:"column:url"`

	addr *sip.Address `gorm:"-"`
}

// 从请求中解析出设备信息
func parserDevicesFromReqeust(req *sip.Request) (Devices, bool) {
	u := Devices{}
	header, ok := req.From()
	if !ok {
		Logger.Warn("not found from header from request", zap.Any("req", req.String()))
		return u, false
	}
	if header.Address == nil {
		Logger.Warn("not found from user from request", zap.Any("req", req.String()))
		return u, false
	}
	if header.Address.User() == nil {
		Logger.Warn("not found from  address user from request", zap.Any("req", req.String()))
		return u, false
	}
	u.DeviceID = header.Address.User().String()
	u.Region = header.Address.Host()
	via, ok := req.ViaHop()
	if !ok {
		Logger.Warn("not found ViaHop from request", zap.Any("req", req.String()))
		return u, false
	}
	u.Host = via.Host
	u.Port = via.Port.String()
	report, ok := via.Params.Get("rport")
	if ok && report != nil {
		u.Rport = report.String()
	}
	raddr, ok := via.Params.Get("received")
	if ok && raddr != nil {
		u.RAddr = raddr.String()
	}

	u.TransPort = via.Transport
	u.URIStr = header.Address.String()
	u.addr = sip.NewAddressFromFromHeader(header)
	u.Source = req.Source().String()
	u.source = req.Source()
	return u, true
}

// 获取设备信息（注册设备）
func sipDeviceInfo(to Devices) {
	hb := sip.NewHeaderBuilder().SetTo(to.addr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE)
	req := sip.NewRequest("", sip.MESSAGE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetDeviceInfoXML(to.DeviceID))
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		Logger.Error("sipDeviceInfo  error,", zap.Error(err))
		return
	}
	_, err = sipResponse(tx)
	if err != nil {
		Logger.Error("sipDeviceInfo  response error", zap.Any("addr", to.addr), zap.Error(err))
		return
	}
}

// sipCatalog 获取注册设备包含的列表
func sipCatalog(to Devices) {

	hb := sip.NewHeaderBuilder().SetTo(to.addr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE)
	req := sip.NewRequest("", sip.MESSAGE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetCatalogXML(to.DeviceID))
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		Logger.Error("sipCatalog  error,", zap.Error(err))
		return
	}
	_, err = sipResponse(tx)
	if err != nil {
		return
	}
}

// MessageDeviceInfoResponse 主设备明细返回结构
type MessageDeviceInfoResponse struct {
	CmdType      string `xml:"CmdType"`
	SN           int    `xml:"SN"`
	DeviceID     string `xml:"DeviceID"`
	DeviceType   string `xml:"DeviceType"`
	Manufacturer string `xml:"Manufacturer"`
	Model        string `xml:"Model"`
	Firmware     string `xml:"Firmware"`
}

func sipMessageDeviceInfo(u Devices, body []byte) error {
	message := &MessageDeviceInfoResponse{}
	if err := utils.XMLDecode([]byte(body), message); err != nil {
		Logger.Error("sipMessageDeviceInfo Unmarshal xml err:", zap.Any("body", body), zap.Error(err))
		return err
	}
	db.UpdateAll(db.DBClient, new(Devices), db.M{"deviceid=?": u.DeviceID}, Devices{
		Model:        message.Model,
		DeviceType:   message.DeviceType,
		Firmware:     message.Firmware,
		Manufacturer: message.Manufacturer,
	})
	return nil
}

type MessageMediaStatusResponse struct {
	CmdType    string `xml:"CmdType"`
	SN         int    `xml:"SN"`
	DeviceID   string `xml:"DeviceID"`
	NotifyType string `xml:"NotifyType"`
}

func sipMessageMediaStatusInfo(u Devices, body []byte) error {
	message := &MessageMediaStatusResponse{}
	if err := utils.XMLDecode([]byte(body), message); err != nil {
		Logger.Error("sipMessageDeviceInfo Unmarshal xml err:", zap.Any("body", body), zap.Error(err))
		return err
	}

	Logger.Info("当前流状态", zap.Any("Device ID", u.DeviceID), zap.Any("Media Statis", message.NotifyType))

	return nil
}

type MessageBroadcastResponse struct {
	CmdType  string `xml:"CmdType"`
	SN       int    `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
	Result   string `xml:"Result"`
}

func sipMessageBroadcastInfo(u Devices, body []byte) error {
	message := &MessageBroadcastResponse{}
	if err := utils.XMLDecode([]byte(body), message); err != nil {
		Logger.Error("sipMessageBroadcastInfo Unmarshal xml err:", zap.Any("body", body), zap.Error(err))
		return err
	}

	Logger.Info("收到广播消息", zap.Any("Device ID", u.DeviceID), zap.Any("Result", message.Result))

	return nil
}

// MessageDeviceListResponse 设备明细列表返回结构
type MessageDeviceListResponse struct {
	XMLName  xml.Name   `xml:"Response"`
	CmdType  string     `xml:"CmdType"`
	SN       int        `xml:"SN"`
	DeviceID string     `xml:"DeviceID"`
	SumNum   int        `xml:"SumNum"`
	Item     []Channels `xml:"DeviceList>Item"`
}

func sipMessageCatalog(u Devices, body []byte) error {
	message := &MessageDeviceListResponse{}
	if err := utils.XMLDecode(body, message); err != nil {
		Logger.Error("Message Unmarshal xml err:", zap.Any("body", body), zap.Error(err))
		return err
	}
	if message.SumNum > 0 {
		for _, d := range message.Item {
			channel := Channels{ChannelID: d.ChannelID, DeviceID: message.DeviceID}
			if err := db.Get(db.DBClient, &channel); err == nil || db.RecordNotFound(err) {
				channel.Active = time.Now().Unix()
				channel.URIStr = fmt.Sprintf("sip:%s@%s", d.ChannelID, _sysinfo.Region)
				channel.Status = transDeviceStatus(d.Status)
				channel.Name = d.Name
				channel.Manufacturer = d.Manufacturer
				channel.Model = d.Model
				channel.Owner = d.Owner
				channel.CivilCode = d.CivilCode
				// Address ip地址
				channel.Address = d.Address
				channel.Parental = d.Parental
				channel.SafetyWay = d.SafetyWay
				channel.RegisterWay = d.RegisterWay
				channel.Secrecy = d.Secrecy
				db.Save(db.DBClient, &channel)
				go notify(notifyChannelsActive(channel.ChannelID))
			} else {
				Logger.Error("deviceid not found,deviceid:", zap.Any("device id", d.DeviceID), zap.Any("pdid:", message.DeviceID), zap.Error(err))
			}
		}
	}
	return nil
}

var deviceStatusMap = map[string]string{
	"ON":     m.DeviceStatusON,
	"OK":     m.DeviceStatusON,
	"ONLINE": m.DeviceStatusON,
	"OFFILE": m.DeviceStatusOFF,
	"OFF":    m.DeviceStatusOFF,
}

func transDeviceStatus(status string) string {
	if v, ok := deviceStatusMap[status]; ok {
		return v
	}
	return status
}
