package api

import (
	gapi "go-sip/api/gateway"
	"go-sip/api/middleware"
	sapi "go-sip/api/s"

	. "go-sip/common"

	"github.com/gin-gonic/gin"
)

// 公共初始化
func Init(r *gin.Engine) {
	r.Use(middleware.CORS())
}

// 服务端api初始化
func ServerApiInit(r *gin.Engine) {
	Init(r)

	// 播放类接口
	{
		r.GET(PlayURL, sapi.Play)
		r.GET(PlaybackURL, sapi.Playback)
		r.GET(PauseURL, sapi.Pause)
		r.GET(ResumeURL, sapi.Resume)
		r.GET(SpeedURL, sapi.Speed)
		r.GET(SeekURL, sapi.Seek)
		r.GET(StopURL, sapi.Stop)
	}
	{
		//设备类
		r.GET(DeviceControl, sapi.DeviceControl)
	}
	// 录像类
	{
		r.GET(RecordsListURL, sapi.RecordsList)
	}
	// server zlm webhook
	{
		r.POST(ZLMWebHookServerURL, sapi.ZLMWebHook)
	}

}

// 网关api初始化
func GatewayApiInit(r *gin.Engine) {
	Init(r)

	// 对外开放接口
	{

		// 视频回放相关接口
		r.GET("/open/ipc/records", gapi.RecordsList)
		r.GET("/open/ipc/playbackSpeed", gapi.PlaybackSpeed)
		r.GET("/open/ipc/playbackSeek", gapi.PlaybackSeek)
		r.GET("/open/ipc/playbackStop", gapi.PlaybackStop)
		r.GET("/open/ipc/playbackPause", gapi.PlaybackPause)
		r.GET("/open/ipc/playbackResume", gapi.PlaybackResume)

		r.GET("/open/ipc/control", gapi.DeviceControl)

		// hook相关接口
		r.POST("/open/zlm/webhook/:method", gapi.ZLMWebHook)

		r.POST("/open/server/register", gapi.RegisterSipServer)
		r.GET("/open/server/getone", gapi.GetSipServer)

	}
}
