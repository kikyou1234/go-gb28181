package sipapi

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	db "go-sip/db/sqlite"
	. "go-sip/logger"
	"go-sip/m"
	sip "go-sip/sip/s"
	"go-sip/utils"

	sdp "github.com/panjjo/gosdp"
	"go.uber.org/zap"
)

// MessageReceive 接收到的请求数据最外层，主要用来判断数据类型
type MessageReceive struct {
	CmdType string `xml:"CmdType"`
	SN      int    `xml:"SN"`
}

func handlerMessage(req *sip.Request, tx *sip.Transaction) {
	u, ok := parserDevicesFromReqeust(req)
	if !ok {
		// 未解析出来源用户返回错误
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
		return
	}
	// 判断是否存在body数据
	if len, have := req.ContentLength(); !have || len.Equals(0) {
		// 不存在就直接返回的成功
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	}
	body := req.Body()
	message := &MessageReceive{}

	if err := utils.XMLDecode(body, message); err != nil {
		Logger.Error("Message Unmarshal xml err:", zap.Error(err))
		// 有些body xml发送过来的不带encoding ，而且格式不是utf8的，导致xml解析失败，此处使用gbk转utf8后再次尝试xml解析
		body, err = utils.GbkToUtf8(body)
		if err != nil {
			Logger.Error("message gbk to utf8 err:", zap.Error(err))
		}
		if err := utils.XMLDecode(body, message); err != nil {
			Logger.Error("Message Unmarshal xml after gbktoutf8 err:", zap.Error(err))
			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
			return
		}
	}
	switch message.CmdType {
	case "Catalog":
		// 设备列表
		sipMessageCatalog(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	case "Keepalive":
		// heardbeat
		if err := sipMessageKeepalive(u, body); err == nil {
			tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
			// 心跳后同步注册设备列表信息
			sipCatalog(u)
			return
		}
	case "RecordInfo":
		// 设备音视频文件列表
		sipMessageRecordInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
	case "DeviceInfo":
		// 主设备信息
		sipMessageDeviceInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	case "MediaStatus":
		// 媒体状态
		sipMessageMediaStatusInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	case "Broadcast":
		sipMessageBroadcastInfo(u, body)
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	default:
		Logger.Info("收到消息", zap.Any("body", string(body)))
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
	}

}

func handlerRegister(req *sip.Request, tx *sip.Transaction) {
	// 判断是否存在授权字段
	if hdrs := req.GetHeaders("Authorization"); len(hdrs) > 0 {
		fromUser, ok := parserDevicesFromReqeust(req)
		if !ok {
			return
		}
		user := Devices{DeviceID: fromUser.DeviceID}
		err := db.Get(db.DBClient, &user)

		if err == nil || db.RecordNotFound(err) {
			if !user.Regist {
				// 如果数据库里用户未激活，替换user数据
				fromUser.ID = user.ID
				fromUser.Name = user.Name
				fromUser.PWD = m.CMConfig.GB28181.Passwd
				user = fromUser
			}
			user.addr = fromUser.addr
			authenticateHeader := hdrs[0].(*sip.GenericHeader)
			auth := sip.AuthFromValue(authenticateHeader.Contents)
			auth.SetPassword(user.PWD)
			auth.SetUsername(user.DeviceID)
			auth.SetMethod(string(req.Method()))
			auth.SetURI(auth.Get("uri"))
			if auth.CalcResponse() == auth.Get("response") {
				// 验证成功
				// 记录活跃设备
				user.source = fromUser.source
				user.addr = fromUser.addr
				_activeDevices.Store(user.DeviceID, user)
				if !user.Regist {
					// 第一次激活，保存数据库
					user.Regist = true
					db.DBClient.Save(&user)
					Logger.Info("new user regist,id:", zap.Any("device  id", user.DeviceID))
				}
				tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
				// 注册成功后查询设备信息，获取制作厂商等信息
				go notify(notifyDevicesRegister(user.DeviceID))
				go sipDeviceInfo(fromUser)
				return
			}
		}
	}
	resp := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
	resp.AppendHeader(&sip.GenericHeader{HeaderName: "WWW-Authenticate", Contents: fmt.Sprintf("Digest nonce=\"%s\", algorithm=MD5, realm=\"%s\",qop=\"auth\"", utils.RandString(32), _sysinfo.Region)})
	tx.Respond(resp)
}

func handlerInvite(req *sip.Request, tx *sip.Transaction) {

	u, ok := parserDevicesFromReqeust(req)
	if !ok {
		// 未解析出来源用户返回错误
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil))
		return
	}
	// 判断是否存在body数据
	if len, have := req.ContentLength(); !have || len.Equals(0) {
		// 不存在就直接返回的成功
		tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", nil))
		return
	}

	Logger.Info("收到invite请求", zap.Any("req", req.Body()))
	go notify(notifyInvite(u.DeviceID))

	ack, err := InviteFunc(u.DeviceID)
	if err != nil {
		Logger.Error("invite 失败", zap.Error(err))
	}

	Logger.Info("收到服务端 invite 反馈 请求", zap.Any("ack", ack))

	resp := sip.NewResponseFromRequest("", req, http.StatusContinue, "Trying", nil)
	tx.Respond(resp)

	var (
		s sdp.Session
		b []byte
	)
	name := "Play"

	protocal := "TCP/RTP/AVP"

	video := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "audio",
			Port:     int(ack.ZlmPort),
			Formats:  []string{"8"},
			Protocol: protocal,
		},
	}
	video.AddAttribute("sendonly")

	video.AddAttribute("setup", "passive")
	video.AddAttribute("connection", "new")

	video.AddAttribute("rtpmap", "8", "PCMA/8000/1")
	rtpIp, _ := net.ResolveIPAddr("ip", ack.ZlmIp)
	// defining message
	msg := &sdp.Message{
		Origin: sdp.Origin{
			Username: _serverDevices.DeviceID, // 媒体服务器id
			Address:  ack.ZlmIp,
		},
		Name: name,
		Connection: sdp.ConnectionData{
			IP:  rtpIp.IP,
			TTL: 0,
		},

		Medias: []sdp.Media{video},
		SSRC:   extractY(string(req.Body())),
	}

	s = msg.Append(s)
	// appending session to byte buffer
	b = s.AppendTo(b)

	res := sip.NewResponse(
		"",
		req.SipVersion(),
		http.StatusOK,
		http.StatusText(http.StatusOK),
		[]sip.Header{},
		[]byte{},
	)

	sip.CopyHeaders("Record-Route", req, res)
	sip.CopyHeaders("Via", req, res)
	sip.CopyHeaders("From", req, res)
	sip.CopyHeaders("To", req, res)
	sip.CopyHeaders("Call-ID", req, res)
	sip.CopyHeaders("CSeq", req, res)
	sip.CopyHeaders("Content-Type", req, res)

	res.SetSource(req.Destination())
	res.SetDestination(req.Source())

	if len(b) > 0 {
		res.SetBody(b, true)
	}

	tx.Respond(res)
}

func handlerBye(req *sip.Request, tx *sip.Transaction) {
	// 处理BYE请求
	Logger.Info("收到BYE请求", zap.Any("req", req.Body()))
	resp := sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), nil)
	tx.Respond(resp)
}

func handlerAck(req *sip.Request, tx *sip.Transaction) {
	// 处理BYE请求
	Logger.Info("收到ACK请求", zap.Any("req", req.Body()))
}

func extractY(sdp string) (y string) {
	// 按行分割 SDP
	lines := strings.Split(sdp, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 提取 y= 的值
		if strings.HasPrefix(line, "y=") {
			y = strings.SplitN(line, "=", 2)[1]
		}
	}
	return
}
