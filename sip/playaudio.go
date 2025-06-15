package sipapi

import (
	"errors"
	. "go-sip/logger"
	sip "go-sip/sip/s"
	"go-sip/utils"
	"net"
	"net/http"

	sdp "github.com/panjjo/gosdp"
)

func SipPlayAudio(stream_id string, zlm_port int, zlm_ip string) error {

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

	var (
		s sdp.Session
		b []byte
	)
	name := "Play"

	protocal := "TCP/RTP/AVP"

	audio := sdp.Media{
		Description: sdp.MediaDescription{
			Type:     "audio",
			Port:     zlm_port,
			Formats:  []string{"8"},
			Protocol: protocal,
		},
	}
	audio.AddAttribute("sendonly")

	audio.AddAttribute("setup", "passive")
	audio.AddAttribute("connection", "new")

	audio.AddAttribute("rtpmap", "8", "PCMA/8000/1")
	rtpIp, _ := net.ResolveIPAddr("ip", zlm_ip)
	// defining message
	msg := &sdp.Message{
		Origin: sdp.Origin{
			Username: _serverDevices.DeviceID, // 媒体服务器id
			Address:  zlm_ip,
		},
		Name: name,
		Connection: sdp.ConnectionData{
			IP:  rtpIp.IP,
			TTL: 0,
		},
		Medias: []sdp.Media{audio},
	}

	s = msg.Append(s)
	// appending session to byte buffer
	b = s.AppendTo(b)

	_serverDevices.addr.Params.Add("tag", sip.String{Str: utils.RandString(20)})
	hb := sip.NewHeaderBuilder().SetTo(to.addr).SetFrom(_serverDevices.addr).AddVia(&sip.ViaHop{
		Params: sip.NewParams().Add("branch", sip.String{Str: sip.GenerateBranch()}),
	}).SetContentType(&sip.ContentTypeSDP).SetMethod(sip.INVITE).SetContact(_serverDevices.addr)
	req := sip.NewRequest("", sip.INVITE, to.addr.URI, sip.DefaultSipVersion, hb.Build(), b)
	req.SetDestination(to.source)
	req.SetRecipient(to.addr.URI)
	tx, err := srv.Request(req)
	if err != nil {
		return err
	}

	response := tx.GetResponse()
	if response.StatusCode() != http.StatusOK {
		return errors.New(response.Reason())
	}

	return err
}
