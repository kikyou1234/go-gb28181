package sipapi

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	db "go-sip/db/sqlite"
	. "go-sip/logger"
	"go-sip/m"
	sip "go-sip/sip/s"
	"go-sip/utils"

	sdp "github.com/panjjo/gosdp"
	"go.uber.org/zap"
)

// sip 请求播放
func SipPlay(data *Streams) (*Streams, error) {
	Logger.Info("SipPlay", zap.Any("channel id", data.ChannelID))
	channel := Channels{ChannelID: data.ChannelID}
	if err := db.Get(db.DBClient, &channel); err != nil {
		if db.RecordNotFound(err) {
			return nil, errors.New("通道不存在")
		}
		return nil, err
	}
	data.DeviceID = channel.DeviceID
	data.StreamType = channel.StreamType
	// 使用通道的播放模式进行处理
	switch channel.StreamType {
	case m.StreamTypePull:
		// 拉流
	default:
		// 推流模式要求设备在线且活跃
		if time.Now().Unix()-channel.Active > 30*60 || channel.Status != m.DeviceStatusON {
			return nil, errors.New("通道已离线")
		}
		user, ok := _activeDevices.Get(channel.DeviceID)
		if !ok {
			return nil, errors.New("设备已离线")
		}
		ssrcLock.Lock()
		data.ssrc = getSSRC(data.T)
		ssrcLock.Unlock()
		var err error
		data, err = sipPlayPush(data, channel, user)
		if err != nil {
			return nil, fmt.Errorf("获取视频失败:%v", err)
		}
	}

	data.Ext = time.Now().Unix() + 2*60 // 2分钟等待时间

	StreamList.Response.Store(data.StreamID, data)
	Logger.Info("StreamList.Response.Store:", zap.Any("channel id", channel.ChannelID), zap.Any("stream id", data.StreamID))
	if data.T == 0 {
		StreamList.Succ.Store(data.ChannelID, data)
	}

	db.Save(db.DBClient, data)
	return data, nil
}

var ssrcLock *sync.Mutex

func sipPlayPush(data *Streams, channel Channels, device Devices) (*Streams, error) {
	Logger.Info("sipPlayPush", zap.Any("channel id", data.ChannelID))
	var (
		s sdp.Session
		b []byte
	)
	name := "Play"
	protocal := "RTP/AVP"
	if data.Mode == 1 {
		protocal = "TCP/RTP/AVP"
	} else {
		protocal = "RTP/AVP"
	}
	if data.T == 1 {
		name = "Playback"
	}
	video := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "video",
			Port:     data.ZlmPort,
			Formats:  []string{"96", "98", "97", "99"},
			Protocol: protocal,
		},
	}
	video.AddAttribute("recvonly")
	if data.Mode == 1 {
		video.AddAttribute("setup", "passive")
		video.AddAttribute("connection", "new")
	}
	video.AddAttribute("rtpmap", "96", "PS/90000")
	video.AddAttribute("rtpmap", "98", "H264/90000")
	video.AddAttribute("rtpmap", "97", "MPEG4/90000")
	video.AddAttribute("rtpmap", "99", "H265/90000")
	if data.Resolution == 0 {
		video.AddAttribute("stream", "1") // 1 表示标清
	} else {
		video.AddAttribute("stream", "0") // 0  表示高清
	}

	rtpIp, _ := net.ResolveIPAddr("ip", data.ZlmIP)
	// defining message
	msg := &sdp.Message{
		Origin: sdp.Origin{
			Username: _serverDevices.DeviceID, // 媒体服务器id
			Address:  data.ZlmIP,
		},
		Name: name,
		Connection: sdp.ConnectionData{
			IP:  rtpIp.IP,
			TTL: 0,
		},
		Timing: []sdp.Timing{
			{
				Start: data.S,
				End:   data.E,
			},
		},
		Medias: []sdp.Media{video},
		SSRC:   data.ssrc,
	}
	if data.T == 1 {
		msg.URI = fmt.Sprintf("%s:1", channel.ChannelID)
	}

	s = msg.Append(s)
	// appending session to byte buffer
	b = s.AppendTo(b)
	uri, _ := sip.ParseURI(channel.URIStr)
	channel.addr = &sip.Address{URI: uri}

	_serverDevices.addr.Params.Add("tag", sip.String{Str: utils.RandString(20)})
	hb := sip.NewHeaderBuilder().SetTo(channel.addr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeSDP).SetMethod(sip.INVITE).SetContact(_serverDevices.addr)
	req := sip.NewRequest("", sip.INVITE, channel.addr.URI, sip.DefaultSipVersion, hb.Build(), b)
	req.SetDestination(device.source)
	req.AppendHeader(&sip.GenericHeader{HeaderName: "Subject", Contents: fmt.Sprintf("%s:%s,%s:%s", channel.ChannelID, data.StreamID, _serverDevices.DeviceID, data.StreamID)})
	req.SetRecipient(channel.addr.URI)
	tx, err := srv.Request(req)
	if err != nil {
		Logger.Error("sipPlayPush fail.id:", zap.Any("channel id", channel.ChannelID), zap.Error(err))
		return data, err
	}
	// response
	response, err := sipResponse(tx)
	if err != nil {
		Logger.Error("sipPlayPush response fail.id ", zap.Any("channel id", channel.ChannelID), zap.Error(err))
		return data, err
	}
	data.Resp = response
	// ACK
	tx.Request(sip.NewRequestFromResponse(sip.ACK, response))

	callid, _ := response.CallID()
	data.CallID = string(*callid)

	cseq, _ := response.CSeq()
	if cseq != nil {
		data.CseqNo = cseq.SeqNo
	}

	from, _ := response.From()
	to, _ := response.To()
	for k, v := range to.Params.Items() {
		data.Ttag[k] = v.String()
	}
	for k, v := range from.Params.Items() {
		data.Ftag[k] = v.String()
	}
	data.Status = 0

	return data, err
}

// sip 停止播放
func SipStopPlay(stream_id string) error {
	Logger.Info("停止播放", zap.Any("stream id", stream_id))
	data, ok := StreamList.Response.Load(stream_id)
	if !ok {
		Logger.Error("停止失败，流不存在")
		return errors.New("流不存在")
	}
	play := data.(*Streams)

	// 推流，需要发送关闭请求
	resp := play.Resp
	u, ok := _activeDevices.Load(play.DeviceID)
	if !ok {
		return errors.New("活跃设备不存在")
	}
	user := u.(Devices)
	req := sip.NewRequestFromResponse(sip.BYE, resp)
	req.SetDestination(user.source)
	tx, err := srv.Request(req)
	if err != nil {
		Logger.Error("sipPlayPush response fail.id ", zap.Error(err))
		return err
	}
	_, err = sipResponse(tx)
	if err != nil {
		Logger.Error("sipStopPlay response fail ", zap.Error(err))
		play.Msg = err.Error()
		return err
	} else {
		play.Status = 1
		play.Stop = true
	}
	db.Save(db.DBClient, play)

	StreamList.Response.Delete(stream_id)
	if play.T == 0 {
		StreamList.Succ.Delete(play.ChannelID)
	}
	return nil
}

// 暂停播放
func SipPausePlay(stream_id string) error {
	Logger.Info("暂停播放", zap.Any("stream id", stream_id))
	data, ok := StreamList.Response.Load(stream_id)
	if !ok {
		Logger.Error("暂停播放失败，流不存在")
		return errors.New("流不存在")
	}
	play := data.(*Streams)

	u, ok := _activeDevices.Load(play.DeviceID)
	if !ok {
		return errors.New("活跃设备不存在")
	}
	to := u.(Devices)
	channelURI, _ := sip.ParseURI(to.URIStr)
	to.addr = &sip.Address{URI: channelURI}

	_url, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%s", to.DeviceID, to.Host, to.Port))
	contact := &sip.Address{
		URI: &_url,
	}

	call_id := sip.CallID(play.CallID)

	hb := sip.NewHeaderBuilder().SetTo(contact).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeRTSP).SetMethod(sip.INFO).SetSeqNo(uint(play.CseqNo)).SetContact(contact).SetCallID(&call_id)
	req := sip.NewRequest("", sip.INFO, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetPausePlayContent())
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}
	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return nil
}

// 恢复播放
func SipResumePlay(stream_id string) error {
	Logger.Info("恢复播放", zap.Any("stream_id", stream_id))
	data, ok := StreamList.Response.Load(stream_id)
	if !ok {
		Logger.Error("恢复播放失败，流不存在")
		return errors.New("流不存在")
	}
	play := data.(*Streams)

	u, ok := _activeDevices.Load(play.DeviceID)
	if !ok {
		return errors.New("活跃设备不存在")
	}
	to := u.(Devices)
	channelURI, _ := sip.ParseURI(to.URIStr)
	to.addr = &sip.Address{URI: channelURI}

	_url, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%s", to.DeviceID, to.Host, to.Port))
	contact := &sip.Address{
		URI: &_url,
	}

	call_id := sip.CallID(play.CallID)

	hb := sip.NewHeaderBuilder().SetTo(contact).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeRTSP).SetMethod(sip.INFO).SetSeqNo(uint(play.CseqNo)).SetContact(contact).SetCallID(&call_id)
	req := sip.NewRequest("", sip.INFO, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetResumePlayContent())
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}
	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return nil
}

// 拖动播放
func SipSeekPlay(stream_id string, sub_time int64) error {
	Logger.Info("拖动播放", zap.Any("stream_id", stream_id), zap.Any("sub_time", sub_time))
	data, ok := StreamList.Response.Load(stream_id)
	if !ok {
		Logger.Error("拖动播放失败，流不存在")
		return errors.New("流不存在")
	}
	play := data.(*Streams)

	u, ok := _activeDevices.Load(play.DeviceID)
	if !ok {
		return errors.New("活跃设备不存在")
	}
	to := u.(Devices)
	channelURI, _ := sip.ParseURI(to.URIStr)
	to.addr = &sip.Address{URI: channelURI}

	_url, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%s", to.DeviceID, to.Host, to.Port))
	contact := &sip.Address{
		URI: &_url,
	}

	call_id := sip.CallID(play.CallID)
	hb := sip.NewHeaderBuilder().SetTo(contact).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeRTSP).SetMethod(sip.INFO).SetSeqNo(uint(play.CseqNo)).SetContact(contact).SetCallID(&call_id)
	req := sip.NewRequest("", sip.INFO, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetSeekPlayContent(sub_time))
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}
	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return nil
}

// 倍速播放
func SipSpeedPlay(stream_id string, speed float64) error {
	Logger.Info("倍速播放", zap.Any("string_id", stream_id), zap.Any("speed", speed))
	data, ok := StreamList.Response.Load(stream_id)
	if !ok {
		Logger.Error("倍速播放失败，流不存在")
		return errors.New("流不存在")
	}
	play := data.(*Streams)

	u, ok := _activeDevices.Load(play.DeviceID)
	if !ok {
		return errors.New("活跃设备不存在")
	}
	to := u.(Devices)
	channelURI, _ := sip.ParseURI(to.URIStr)
	to.addr = &sip.Address{URI: channelURI}

	_url, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%s", to.DeviceID, to.Host, to.Port))
	contact := &sip.Address{
		URI: &_url,
	}

	call_id := sip.CallID(play.CallID)

	hb := sip.NewHeaderBuilder().SetTo(contact).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeRTSP).SetMethod(sip.INFO).SetSeqNo(uint(play.CseqNo)).SetContact(contact).SetCallID(&call_id)
	req := sip.NewRequest("", sip.INFO, to.addr.URI, sip.DefaultSipVersion, hb.Build(), sip.GetSpeedPlayContent(speed))
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}
	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return nil
}

func SipIpcBroadCast(channal_id string) error {

	to, ok := _activeDevices.Get(channal_id)
	if !ok {
		return errors.New("设备已离线")
	}

	channelURI, _ := sip.ParseURI(to.URIStr)
	to.addr = &sip.Address{URI: channelURI}

	_url, _ := sip.ParseSipURI(fmt.Sprintf("sip:%s@%s:%s", to.DeviceID, to.Host, to.Port))
	contact := &sip.Address{
		URI: &_url,
	}

	hb := sip.NewHeaderBuilder().SetTo(contact).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeXML).SetMethod(sip.MESSAGE).SetContact(contact)
	req := sip.NewRequest("", sip.MESSAGE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), []byte(sip.GenerateBroadcastXML(_serverDevices.DeviceID, channal_id)))
	req.SetDestination(to.source)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}
	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return nil

}
