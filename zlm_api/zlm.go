package zlm_api

import (
	. "go-sip/logger"
	"go-sip/utils"

	"fmt"

	"go.uber.org/zap"
)

type ZlmStartSendRtpReq struct {
	Vhost    string `json:"vhost"`
	StreamID string `json:"stream_id"`
	App      string `json:"app"`
	DstUrl   string `json:"dst_url"`
	DstPort  string `json:"dst_port"`
	Ssrc     string `json:"ssrc"`
	IsUdp    string `json:"is_udp"` // 1:udp active模式, 0:tcp active模式
}

type ZlmGetRtpInfoReq struct {
	Vhost    string `json:"vhost"`
	StreamID string `json:"stream_id"`
	App      string `json:"app"`
}

type ZlmGetMediaListReq struct {
	Vhost    string `json:"vhost"`
	Schema   string `json:"schema"`
	StreamID string `json:"stream_id"`
	App      string `json:"app"`
}

type ZlmStartSendRtpResp struct {
	Code      int `json:"code"`
	LocalPort int `json:"local_port"`
}

type ZlmGetMediaListResp struct {
	Code int                       `json:"code"`
	Data []ZlmGetMediaListDataResp `json:"data"`
}

type ZlmGetRtpInfoResp struct {
	Code  int  `json:"code"`
	Exist bool `json:"exist"`
}
type ZlmGetMediaListDataResp struct {
	App        string                  `json:"app"`
	Stream     string                  `json:"stream"`
	Schema     string                  `json:"schema"`
	OriginType int                     `json:"originType"`
	Tracks     []ZlmGetMediaListTracks `json:"tracks"`
}
type ZlmGetMediaListTracks struct {
	Type    int `json:"codec_type"`
	CodecID int `json:"codec_id"`
	Height  int `json:"height"`
	Width   int `json:"width"`
	FPS     int `json:"fps"`
}

// Zlm 开始active模式发送rtp
// 作为zlm客户端，启动ps-rtp推流，支持rtp/udp方式；该接口支持rtsp/rtmp等协议转ps-rtp推流。第一次推流失败会直接返回错误，成功一次后，后续失败也将无限重试。
func ZlmStartSendRtp(url, secret string, req ZlmStartSendRtpReq) ZlmStartSendRtpResp {
	res := ZlmStartSendRtpResp{}
	reqStr := "/index/api/startSendRtp?secret=" + secret
	if req.StreamID != "" {
		reqStr += "&stream=" + req.StreamID
	}
	if req.App != "" {
		reqStr += "&app=" + req.App
	}
	if req.Vhost != "" {
		reqStr += "&vhost=" + req.Vhost
	}
	if req.DstUrl != "" {
		reqStr += "&dst_url=" + req.DstUrl
	}
	if req.DstPort != "" {
		reqStr += "&dst_port=" + req.DstPort
	}
	if req.IsUdp != "" {
		reqStr += "&is_udp=" + req.IsUdp
	}
	if req.Ssrc != "" {
		reqStr += "&ssrc=" + req.Ssrc
	}
	body, err := utils.GetRequest(url + reqStr)
	if err != nil {
		Logger.Error("ZlmStartSendRtp fail, 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("ZlmStartSendRtp fail, 2", zap.Error(err))
		return res
	}
	return res
}

// Zlm 获取RTP流信息
func ZlmGetRtpInfo(url, secret string, req ZlmGetRtpInfoReq) ZlmGetRtpInfoResp {
	res := ZlmGetRtpInfoResp{}
	reqStr := "/index/api/getRtpInfo?secret=" + secret
	if req.StreamID != "" {
		reqStr += "&stream_id=" + req.StreamID
	}
	if req.App != "" {
		reqStr += "&app=" + req.App
	}
	if req.Vhost != "" {
		reqStr += "&vhost=" + req.Vhost
	}
	body, err := utils.GetRequest(url + reqStr)
	if err != nil {
		Logger.Error("get stream rtpInfo fail, 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("get stream rtpInfo fail, 2", zap.Error(err))
		return res
	}
	return res
}

// Zlm 获取流列表信息
func ZlmGetMediaList(url, secret string, req ZlmGetMediaListReq) ZlmGetMediaListResp {
	res := ZlmGetMediaListResp{}
	reqStr := "/index/api/getMediaList?secret=" + secret
	if req.StreamID != "" {
		reqStr += "&stream=" + req.StreamID
	}
	if req.App != "" {
		reqStr += "&app=" + req.App
	}
	if req.Schema != "" {
		reqStr += "&schema=" + req.Schema
	}
	if req.Vhost != "" {
		reqStr += "&vhost=" + req.Vhost
	}
	body, err := utils.GetRequest(url + reqStr)
	if err != nil {
		Logger.Error("get stream mediaList fail, 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("get stream mediaList fail, 2", zap.Error(err))
		return res
	}
	return res
}

var ZlmDeviceVFMap = map[int]string{
	0: "H264",
	1: "H265",
	2: "ACC",
	3: "G711A",
	4: "G711U",
}

func TransZlmDeviceVF(t int) string {
	if v, ok := ZlmDeviceVFMap[t]; ok {
		return v
	}
	return "undefind"
}

type RtpInfo struct {
	Code  int  `json:"code"`
	Exist bool `json:"exist"`
}

// 获取流在Zlm上的信息
func ZlmGetMediaInfo(url, secret, stream_id string) RtpInfo {
	res := RtpInfo{}
	body, err := utils.GetRequest(url + "/index/api/getRtpInfo?secret=" + secret + "&stream_id=" + stream_id)
	if err != nil {
		Logger.Error("get stream rtpInfo fail 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("get stream rtpInfo fail 2", zap.Error(err))
		return res
	}
	return res
}

// Zlm 关闭流
func ZlmCloseStream(url, secret, stream_id string) {
	utils.GetRequest(url + "/index/api/close_streams?secret=" + secret + "&stream=" + stream_id)
}

type OpenRtpRsp struct {
	Code int `json:"code"`
	Port int `json:"port"`
}

// /openRtpServer
func ZlmOpenRtpServer(url, secret, stream_id string, tcp_mode int) OpenRtpRsp {
	res := OpenRtpRsp{}
	// port: 接收端口，0则为随机端口
	// tcp_mode: 0 udp 模式，1 tcp 被动模式, 2 tcp 主动模式
	body, err := utils.GetRequest(url + "/index/api/openRtpServer?secret=" + secret + "&stream_id=" + stream_id + "&port=0" + "&tcp_mode=" + fmt.Sprintf("%d", tcp_mode))
	if err != nil {
		Logger.Error("open server rtp fail 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("open server rtp fail 2", zap.Error(err))
		return res
	}
	return res
}

type CloseRtpRsp struct {
	Code int `json:"code"`
	Hit  int `json:"hit"`
}

// /openRtpServer
func ZlmCloseRtpServer(url, secret, stream_id string) CloseRtpRsp {
	res := CloseRtpRsp{}

	body, err := utils.GetRequest(url + "/index/api/closeRtpServer?secret=" + secret + "&stream_id=" + stream_id)
	if err != nil {
		Logger.Error("close server rtp fail 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("close server rtp fail 2", zap.Error(err))
		return res
	}
	return res
}

type PauseRtpRsp struct {
	Code int `json:"code"`
}

func ZlmPauseRtpCheck(url, secret, stream_id string) PauseRtpRsp {
	res := PauseRtpRsp{}

	body, err := utils.GetRequest(url + "/index/api/pauseRtpCheck?secret=" + secret + "&app=rtp" + "&stream_id=" + stream_id)
	if err != nil {
		Logger.Error("pause server rtp check fail 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("pause server rtp check fail 2", zap.Error(err))
		return res
	}
	return res
}

func ZlmResumeRtpCheck(url, secret, stream_id string) PauseRtpRsp {
	res := PauseRtpRsp{}

	body, err := utils.GetRequest(url + "/index/api/resumeRtpCheck?secret=" + secret + "&app=rtp" + "&stream_id=" + stream_id)
	if err != nil {
		Logger.Error("resume server rtp check fail 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("resume server rtp fail 2", zap.Error(err))
		return res
	}
	return res
}

type OpenSendRtpRsp struct {
	Code      int `json:"code"`
	LocalPort int `json:"local_port"`
}

func ZlmStartSendRtpPassive(ip, port, secret, stream_id string) OpenSendRtpRsp {
	res := OpenSendRtpRsp{}
	url := "http://" + ip + ":" + port + "/index/api/startSendRtpPassive?secret=" + secret + "&stream=" + stream_id + "&ssrc=1&app=broadcast&vhost=__defaultVhost__&only_audio=1&pt=8&use_ps=0&is_udp=0"
	body, err := utils.GetRequest(url)
	if err != nil {
		Logger.Error("start rtp passive fail 1", zap.Error(err))
		return res
	}
	if err = utils.JSONDecode(body, &res); err != nil {
		Logger.Error("start rtp passive fail 2", zap.Error(err))
		return res
	}
	return res
}

// Zlm 开始录制视频流
// func ZlmStartRecord(values url.Values) error {
// 	body, err := utils.GetRequest(config.Media.RESTFUL + "/index/api/startRecord?" + values.Encode())
// 	if err != nil {
// 		return err
// 	}
// 	tmp := map[string]interface{}{}
// 	err = utils.JSONDecode(body, &tmp)
// 	if err != nil {
// 		return err
// 	}
// 	if code, ok := tmp["code"]; !ok || fmt.Sprint(code) != "0" {
// 		return utils.NewError(nil, tmp)
// 	}
// 	return nil
// }

// // Zlm 停止录制
// func ZlmStopRecord(values url.Values) error {
// 	body, err := utils.GetRequest(config.Media.RESTFUL + "/index/api/stopRecord?" + values.Encode())
// 	if err != nil {
// 		return err
// 	}
// 	tmp := map[string]interface{}{}
// 	err = utils.JSONDecode(body, &tmp)
// 	if err != nil {
// 		return err
// 	}
// 	if code, ok := tmp["code"]; !ok || fmt.Sprint(code) != "0" {
// 		return utils.NewError(nil, tmp)
// 	}
// 	return nil
// }
