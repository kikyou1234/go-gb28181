package sipapi

import (
	"fmt"
	db "go-sip/db/sqlite"
	. "go-sip/logger"
	"go-sip/m"
	sip "go-sip/sip/s"
	"go-sip/utils"
	"net/http"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

var _activeDevices ActiveDevices

// 系统运行信息
var _sysinfo *m.SysInfo
var config *m.C_Config

func Start() {
	// 数据库表初始化 启动时自动同步数据结构到数据库
	db.DBClient.AutoMigrate(new(Devices))
	db.DBClient.AutoMigrate(new(Channels))
	db.DBClient.AutoMigrate(new(Streams))

	LoadSYSInfo()
	srv = sip.NewServer()
	srv.RegistHandler(sip.REGISTER, handlerRegister)
	srv.RegistHandler(sip.MESSAGE, handlerMessage)
	srv.RegistHandler(sip.INVITE, handlerInvite)
	srv.RegistHandler(sip.BYE, handlerBye)
	srv.RegistHandler(sip.ACK, handlerAck)

	Logger.Info("sip server start", zap.String("region", _sysinfo.Region), zap.String("lid", _sysinfo.LID))
	go srv.ListenUDPServer(config.UDP)
}

// ActiveDevices 记录当前活跃设备，请求播放时设备必须处于活跃状态
type ActiveDevices struct {
	sync.Map
}

// Get Get
func (a *ActiveDevices) Get(key string) (Devices, bool) {
	if v, ok := a.Load(key); ok {
		return v.(Devices), ok
	}
	return Devices{}, false
}

func LoadSYSInfo() {

	config = m.CMConfig
	_activeDevices = ActiveDevices{sync.Map{}}

	StreamList = streamsList{&sync.Map{}, &sync.Map{}, 0}
	ssrcLock = &sync.Mutex{}
	_recordList = &sync.Map{}

	// init sysinfo
	_sysinfo = &m.SysInfo{}
	_sysinfo = m.DefaultInfo()
	uri, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s", _sysinfo.LID, _sysinfo.Region))
	_serverDevices = Devices{
		DeviceID: _sysinfo.LID,
		Region:   _sysinfo.Region,
		addr: &sip.Address{
			DisplayName: sip.String{Str: "sipserver"},
			URI:         &uri,
			Params:      sip.NewParams(),
		},
	}

}

// zlm接收到的ssrc为16进制。发起请求的ssrc为10进制
func ssrc2stream(ssrc string) string {
	if ssrc[0:1] == "0" {
		ssrc = ssrc[1:]
	}
	num, _ := strconv.Atoi(ssrc)
	return fmt.Sprintf("%08X", num)
}

func sipResponse(tx *sip.Transaction) (*sip.Response, error) {
	response := tx.GetResponse()
	if response == nil {
		return nil, utils.NewError(nil, "response timeout", "tx key:", tx.Key())
	}
	if response.StatusCode() != http.StatusOK {
		return response, utils.NewError(nil, "response fail", response.StatusCode(), response.Reason(), "tx key:", tx.Key())
	}
	return response, nil
}
