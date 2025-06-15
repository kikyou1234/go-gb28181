package model

// type ZlmServerStartDate struct {
// 	API           APIConfig       `json:"api"`
// 	Cluster       ClusterConfig   `json:"cluster"`
// 	FFmpeg        FFmpegConfig    `json:"ffmpeg"`
// 	General       GeneralConfig   `json:"general"`
// 	Hls           HlsConfig       `json:"hls"`
// 	Hook          HookConfig      `json:"hook"`
// 	HookIndex     int             `gorm:"column:hook_index" json:"hook_index"`
// 	HTTP          HTTPConfig      `json:"http"`
// 	MediaServerId string          `gorm:"column:media_server_id" json:"mediaServerId"`
// 	Multicast     MulticastConfig `json:"multicast"`
// 	Protocol      ProtocolConfig  `json:"protocol"`
// 	Record        RecordConfig    `json:"record"`
// 	RTC           RTCConfig       `json:"rtc"`
// 	RTMP          RTMPConfig      `json:"rtmp"`
// 	RTP           RTPConfig       `json:"rtp"`
// 	RTSP          RTSPConfig      `json:"rtsp"`
// 	Shell         ShellConfig     `json:"shell"`
// 	Srt           SrtConfig       `json:"srt"`
// }

type ZlmServerStartDate struct {
	APIDebug     string `gorm:"column:api_debug" json:"api.apiDebug"`
	APISecret    string `gorm:"column:api_secret" json:"api.secret"`
	DefaultSnap  string `gorm:"column:api_default_snap" json:"api.defaultSnap"`
	DownloadRoot string `gorm:"column:api_download_root" json:"api.downloadRoot"`
	SnapRoot     string `gorm:"column:api_snap_root" json:"api.snapRoot"`

	OriginURL         string `gorm:"column:origin_url" json:"cluster.origin_url"`
	RetryCount        string `gorm:"column:retry_count" json:"cluster.retry_count"`
	ClusterTimeoutSec string `gorm:"column:timeout_sec" json:"cluster.timeout_sec"`

	Bin        string `gorm:"column:ffmpeg_bin" json:"ffmpeg.bin"`
	Cmd        string `gorm:"column:ffmpeg_cmd" json:"ffmpeg.cmd"`
	Log        string `gorm:"column:ffmpeg_log" json:"ffmpeg.log"`
	RestartSec string `gorm:"column:ffmpeg_restart_sec" json:"ffmpeg.restart_sec"`
	Snap       string `gorm:"column:ffmpeg_snap" json:"ffmpeg.snap"`

	GeneralMediaServerId    string `gorm:"column:media_server_id" json:"general.mediaServerId"`
	GeneralAddMuteAudio     string `gorm:"column:add_mute_audio" json:"general.addMuteAudio"`
	EnableVhost             string `gorm:"column:enable_vhost" json:"general.enableVhost"`
	FlowThreshold           string `gorm:"column:flow_threshold" json:"general.flowThreshold"`
	MaxStreamWaitMS         string `gorm:"column:max_stream_wait_ms" json:"general.maxStreamWaitMS"`
	ResetWhenRePlay         string `gorm:"column:reset_when_re_play" json:"general.resetWhenRePlay"`
	StreamNoneReaderDelayMS string `gorm:"column:stream_none_reader_delay_ms" json:"general.streamNoneReaderDelayMS"`
	UnreadyFrameCache       string `gorm:"column:unready_frame_cache" json:"general.unready_frame_cache"`
	WaitAddTrackMs          string `gorm:"column:wait_add_track_ms" json:"general.wait_add_track_ms"`
	WaitAudioTrackDataMs    string `gorm:"column:wait_audio_track_data_ms" json:"general.wait_audio_track_data_ms"`
	WaitTrackReadyMs        string `gorm:"column:wait_track_ready_ms" json:"general.wait_track_ready_ms"`
	BroadcastPlayerCount    string `gorm:"column:broadcast_player_count_changed" json:"general.broadcast_player_count_changed"`
	CheckNvidiaDev          string `gorm:"column:check_nvidia_dev" json:"general.check_nvidia_dev"`
	MergeWriteMS            string `gorm:"column:merge_write_ms" json:"general.mergeWriteMS"`
	ListenIP                string `gorm:"column:listen_ip" json:"general.listen_ip"`

	BroadcastRecordTs string `gorm:"column:hls_broadcast_record_ts" json:"hls.broadcastRecordTs"`
	DeleteDelaySec    string `gorm:"column:hls_delete_delay_sec" json:"hls.deleteDelaySec"`
	FastRegister      string `gorm:"column:hls_fast_register" json:"hls.fastRegister"`
	HLSFileBufSize    string `gorm:"column:hls_file_buf_size" json:"hls.fileBufSize"`
	SegDelay          string `gorm:"column:hls_seg_delay" json:"hls.segDelay"`
	SegDur            string `gorm:"column:hls_seg_dur" json:"hls.segDur"`
	SegKeep           string `gorm:"column:hls_seg_keep" json:"hls.segKeep"`
	SegNum            string `gorm:"column:hls_seg_num" json:"hls.segNum"`
	SegRetain         string `gorm:"column:hls_seg_retain" json:"hls.segRetain"`

	Alivestringerval     string `gorm:"column:alive_stringerval" json:"hook.alive_stringerval"`
	Enable               string `gorm:"column:enable" json:"hook.enable"`
	OnFlowReport         string `gorm:"column:on_flow_report" json:"hook.on_flow_report"`
	OnHttpAccess         string `gorm:"column:on_http_access" json:"hook.on_http_access"`
	OnPlay               string `gorm:"column:on_play" json:"hook.on_play"`
	OnPublish            string `gorm:"column:on_publish" json:"hook.on_publish"`
	OnRecordMp4          string `gorm:"column:on_record_mp4" json:"hook.on_record_mp4"`
	OnRecordTs           string `gorm:"column:on_record_ts" json:"hook.on_record_ts"`
	OnRtpServerTimeout   string `gorm:"column:on_rtp_server_timeout" json:"hook.on_rtp_server_timeout"`
	OnRtspAuth           string `gorm:"column:on_rtsp_auth" json:"hook.on_rtsp_auth"`
	OnRtspRealm          string `gorm:"column:on_rtsp_realm" json:"hook.on_rtsp_realm"`
	OnSendRtpStopped     string `gorm:"column:on_send_rtp_stopped" json:"hook.on_send_rtp_stopped"`
	OnServerExited       string `gorm:"column:on_server_exited" json:"hook.on_server_exited"`
	OnServerKeepalive    string `gorm:"column:on_server_keepalive" json:"hook.on_server_keepalive"`
	OnServerStarted      string `gorm:"column:on_server_started" json:"hook.on_server_started"`
	OnShellLogin         string `gorm:"column:on_shell_login" json:"hook.on_shell_login"`
	OnStreamChanged      string `gorm:"column:on_stream_changed" json:"hook.on_stream_changed"`
	OnStreamNoneReader   string `gorm:"column:on_stream_none_reader" json:"hook.on_stream_none_reader"`
	OnStreamNotFound     string `gorm:"column:on_stream_not_found" json:"hook.on_stream_not_found"`
	Retry                string `gorm:"column:retry" json:"hook.retry"`
	RetryDelay           string `gorm:"column:retry_delay" json:"hook.retry_delay"`
	StreamChangedSchemas string `gorm:"column:stream_changed_schemas" json:"hook.stream_changed_schemas"`
	HookTimeoutSec       string `gorm:"column:timeout_sec" json:"hook.timeoutSec"`

	HookIndex int `gorm:"column:hook_index" json:"hook_index"`

	AllowCrossDomains   string `gorm:"column:allow_cross_domains" json:"http.allow_cross_domains"`
	AllowIPRange        string `gorm:"column:allow_ip_range" json:"http.allow_ip_range"`
	CharSet             string `gorm:"column:char_set" json:"http.charSet"`
	DirMenu             string `gorm:"column:dir_menu" json:"http.dirMenu"`
	ForbidCacheSuffix   string `gorm:"column:forbid_cache_suffix" json:"http.forbidCacheSuffix"`
	ForwardedIPHeader   string `gorm:"column:forwarded_ip_header" json:"http.forwarded_ip_header"`
	HTTPKeepAliveSecond string `gorm:"column:keep_alive_second" json:"http.keepAliveSecond"`
	MaxReqSize          string `gorm:"column:max_req_size" json:"http.maxReqSize"`
	NotFound            string `gorm:"column:not_found" json:"http.notFound"`
	HTTPPort            string `gorm:"column:port" json:"http.port"`
	RootPath            string `gorm:"column:root_path" json:"http.rootPath"`
	SendBufSize         string `gorm:"column:send_buf_size" json:"http.sendBufSize"`
	HTTPSSLPort         string `gorm:"column:ssl_port" json:"http.sslport"`
	VirtualPath         string `gorm:"column:virtual_path" json:"http.virtualPath"`

	MediaServerId string `gorm:"column:media_server_id" json:"mediaServerId"`

	AddrMax string `gorm:"column:multicast_addr_max" json:"multicast.addrMax"`
	AddrMin string `gorm:"column:multicast_addr_min" json:"multicast.addrMin"`
	UdpTTL  string `gorm:"column:multicast_udp_ttl" json:"multicast.udpTTL"`

	ProtocolAddMuteAudio string `json:"protocol.add_mute_audio" gorm:"column:add_mute_audio"`
	AutoClose            string `json:"protocol.auto_close" gorm:"column:auto_close"`
	ContinuePushMS       string `json:"protocol.continue_push_ms" gorm:"column:continue_push_ms"`
	EnableAudio          string `json:"protocol.enable_audio" gorm:"column:enable_audio"`
	EnableFMP4           string `json:"protocol.enable_fmp4" gorm:"column:enable_fmp4"`
	EnableHLS            string `json:"protocol.enable_hls" gorm:"column:enable_hls"`
	EnableHLSFMP4        string `json:"protocol.enable_hls_fmp4" gorm:"column:enable_hls_fmp4"`
	EnableMP4            string `json:"protocol.enable_mp4" gorm:"column:enable_mp4"`
	EnableRTMP           string `json:"protocol.enable_rtmp" gorm:"column:enable_rtmp"`
	EnableRTSP           string `json:"protocol.enable_rtsp" gorm:"column:enable_rtsp"`
	EnableTS             string `json:"protocol.enable_ts" gorm:"column:enable_ts"`
	Fmp4Demand           string `json:"protocol.fmp4_demand" gorm:"column:fmp4_demand"`
	HlsDemand            string `json:"protocol.hls_demand" gorm:"column:hls_demand"`
	HlsSavePath          string `json:"protocol.hls_save_path" gorm:"column:hls_save_path"`
	ModifyStamp          string `json:"protocol.modify_stamp" gorm:"column:modify_stamp"`
	Mp4AsPlayer          string `json:"protocol.mp4_as_player" gorm:"column:mp4_as_player"`
	Mp4MaxSecond         string `json:"protocol.mp4_max_second" gorm:"column:mp4_max_second"`
	Mp4SavePath          string `json:"protocol.mp4_save_path" gorm:"column:mp4_save_path"`
	PacedSenderMS        string `json:"protocol.paced_sender_ms" gorm:"column:paced_sender_ms"`
	RtmpDemand           string `json:"protocol.rtmp_demand" gorm:"column:rtmp_demand"`
	RtspDemand           string `json:"protocol.rtsp_demand" gorm:"column:rtsp_demand"`
	TsDemand             string `json:"protocol.ts_demand" gorm:"column:ts_demand"`

	AppName          string `gorm:"column:record_app_name" json:"record.appName"`
	EnableFmp4       string `gorm:"column:record_enable_fmp4" json:"record.enableFmp4"`
	FastStart        string `gorm:"column:record_fast_start" json:"record.fastStart"`
	RecodFileBufSize string `gorm:"column:record_file_buf_size" json:"record.fileBufSize"`
	FileRepeat       string `gorm:"column:record_file_repeat" json:"record.fileRepeat"`
	SampleMS         string `gorm:"column:record_sample_ms" json:"record.sampleMS"`

	BFilter              string `json:"rtc.bfilter" gorm:"column:rtc_bfilter"`
	DataChannelEcho      string `json:"rtc.datachannel_echo" gorm:"column:rtc_datachannel_echo"`
	ExternIP             string `json:"rtc.externIP" gorm:"column:rtc_extern_ip"`
	MaxRtpCacheMS        string `json:"rtc.maxRtpCacheMS" gorm:"column:rtc_max_rtp_cache_ms"`
	MaxRtpCacheSize      string `json:"rtc.maxRtpCacheSize" gorm:"column:rtc_max_rtp_cache_size"`
	MaxBitrate           string `json:"rtc.max_bitrate" gorm:"column:rtc_max_bitrate"`
	MinBitrate           string `json:"rtc.min_bitrate" gorm:"column:rtc_min_bitrate"`
	NackstringervalRatio string `json:"rtc.nackstringervalRatio" gorm:"column:rtc_nack_stringerval_ratio"`
	NackMaxCount         string `json:"rtc.nackMaxCount" gorm:"column:rtc_nack_max_count"`
	NackMaxMS            string `json:"rtc.nackMaxMS" gorm:"column:rtc_nack_max_ms"`
	NackMaxSize          string `json:"rtc.nackMaxSize" gorm:"column:rtc_nack_max_size"`
	NackRtpSize          string `json:"rtc.nackRtpSize" gorm:"column:rtc_nack_rtp_size"`
	RTCPort              string `json:"rtc.port" gorm:"column:rtc_port"`
	PreferredCodecA      string `json:"rtc.preferredCodecA" gorm:"column:rtc_preferred_codec_a"`
	PreferredCodecV      string `json:"rtc.preferredCodecV" gorm:"column:rtc_preferred_codec_v"`
	RembBitRate          string `json:"rtc.rembBitRate" gorm:"column:rtc_remb_bitrate"`
	StartBitrate         string `json:"rtc.start_bitrate" gorm:"column:rtc_start_bitrate"`
	TCPPort              string `json:"rtc.tcpPort" gorm:"column:rtc_tcp_port"`
	RTCTimeoutSec        string `json:"rtc.timeoutSec" gorm:"column:rtc_timeout_sec"`

	RTMPDirectProxy     string `gorm:"column:rtmp_direct_proxy" json:"rtmp.directProxy"`
	Enhanced            string `gorm:"column:rtmp_enhanced" json:"rtmp.enhanced"`
	RTMPHandshakeSecond string `gorm:"column:rtmp_handshake_second" json:"rtmp.handshakeSecond"`
	RTMPKeepAliveSecond string `gorm:"column:rtmp_keep_alive_second" json:"rtmp.keepAliveSecond"`
	RTMPPort            string `gorm:"column:rtmp_port" json:"rtmp.port"`
	RTMPSSLPort         string `gorm:"column:rtmp_ssl_port" json:"rtmp.sslport"`

	AudioMtuSize           string `json:"rtp.audioMtuSize" gorm:"column:rtp_audio_mtu_size"`
	H264StapA              string `json:"rtp.h264_stap_a" gorm:"column:rtp_h264_stap_a"`
	RtpLowLatency          string `json:"rtp.lowLatency" gorm:"column:rtp_low_latency"`
	RtpMaxSize             string `json:"rtp.rtpMaxSize" gorm:"column:rtp_rtp_max_size"`
	VideoMtuSize           string `json:"rtp.videoMtuSize" gorm:"column:rtp_video_mtu_size"`
	RtpProxyDumpDir        string `json:"rtp.dumpDir" gorm:"column:rtp_proxy_dump_dir"`
	RtpProxyGopCache       string `json:"rtp.gop_cache" gorm:"column:rtp_proxy_gop_cache"`
	RtpProxyH264PT         string `json:"rtp.h264_pt" gorm:"column:rtp_proxy_h264_pt"`
	RtpProxyH265PT         string `json:"rtp.h265_pt" gorm:"column:rtp_proxy_h265_pt"`
	RtpProxyOpusPT         string `json:"rtp.opus_pt" gorm:"column:rtp_proxy_opus_pt"`
	RtpProxyPort           string `json:"rtp.port" gorm:"column:rtp_proxy_port"`
	RtpProxyPortRange      string `json:"rtp.port_range" gorm:"column:rtp_proxy_port_range"`
	RtpProxyPsPT           string `json:"rtp.ps_pt" gorm:"column:rtp_proxy_ps_pt"`
	RtpProxyRtpG711DurMs   string `json:"rtp.rtp_g711_dur_ms" gorm:"column:rtp_proxy_rtp_g711_dur_ms"`
	RtpProxyTimeoutSec     string `json:"rtp.timeoutSec" gorm:"column:rtp_proxy_timeout_sec"`
	RtpProxyUdpRecvSockBuf string `json:"rtp.udp_recv_socket_buffer" gorm:"column:rtp_proxy_udp_recv_socket_buffer"`

	AuthBasic            string `gorm:"column:rtsp_auth_basic" json:"rtsp.authBasic"`
	RTSPDirectProxy      string `gorm:"column:rtsp_direct_proxy" json:"rtsp.directProxy"`
	RTSPHandshakeSecond  string `gorm:"column:rtsp_handshake_second" json:"rtsp.handshakeSecond"`
	RTSPKeepAliveSecond  string `gorm:"column:rtsp_keep_alive_second" json:"rtsp.keepAliveSecond"`
	LowLatency           string `gorm:"column:rtsp_low_latency" json:"rtsp.lowLatency"`
	RTSPPort             string `gorm:"column:rtsp_port" json:"rtsp.port"`
	RTSPRtpTransportType string `gorm:"column:rtsp_rtp_transport_type" json:"rtsp.rtpTransportType"`
	RTSPSSLPort          string `gorm:"column:rtsp_ssl_port" json:"rtsp.sslport"`

	ShellMaxReqSize string `gorm:"column:shell_max_req_size" json:"shell.maxReqSize"`
	ShellPort       string `gorm:"column:shell_port" json:"shell.port"`

	LatencyMul    string `json:"srt.latencyMul" gorm:"column:srt_latency_mul"`
	PassPhrase    string `json:"srt.passPhrase" gorm:"column:srt_pass_phrase"`
	PktBufSize    string `json:"srt.pktBufSize" gorm:"column:srt_pkt_buf_size"`
	SRTPort       string `json:"srt.port" gorm:"column:srt_port"`
	SRTTimeoutSec string `json:"srt.timeoutSec" gorm:"column:srt_timeout_sec"`
}

type APIConfig struct {
	APIDebug     string `gorm:"column:api_debug" json:"api.apiDebug"`
	Secret       string `gorm:"column:api_secret" json:"api.secret"`
	DefaultSnap  string `gorm:"column:api_default_snap" json:"api.defaultSnap"`
	DownloadRoot string `gorm:"column:api_download_root" json:"api.downloadRoot"`
	SnapRoot     string `gorm:"column:api_snap_root" json:"api.snapRoot"`
}

type ClusterConfig struct {
	OriginURL  string `gorm:"column:origin_url" json:"cluster.origin_url"`
	RetryCount string `gorm:"column:retry_count" json:"cluster.retry_count"`
	TimeoutSec string `gorm:"column:timeout_sec" json:"cluster.timeout_sec"`
}

type FFmpegConfig struct {
	Bin        string `gorm:"column:ffmpeg_bin" json:"bin"`
	Cmd        string `gorm:"column:ffmpeg_cmd" json:"cmd"`
	Log        string `gorm:"column:ffmpeg_log" json:"log"`
	RestartSec string `gorm:"column:ffmpeg_restart_sec" json:"restart_sec"`
	Snap       string `gorm:"column:ffmpeg_snap" json:"snap"`
}

type GeneralConfig struct {
	MediaServerID           string `gorm:"column:media_server_id" json:"mediaServerId"`
	AddMuteAudio            string `gorm:"column:add_mute_audio" json:"addMuteAudio"`
	EnableVhost             string `gorm:"column:enable_vhost" json:"enableVhost"`
	FlowThreshold           string `gorm:"column:flow_threshold" json:"flowThreshold"`
	MaxStreamWaitMS         string `gorm:"column:max_stream_wait_ms" json:"maxStreamWaitMS"`
	ResetWhenRePlay         string `gorm:"column:reset_when_re_play" json:"resetWhenRePlay"`
	StreamNoneReaderDelayMS string `gorm:"column:stream_none_reader_delay_ms" json:"streamNoneReaderDelayMS"`
	UnreadyFrameCache       string `gorm:"column:unready_frame_cache" json:"unready_frame_cache"`
	WaitAddTrackMs          string `gorm:"column:wait_add_track_ms" json:"wait_add_track_ms"`
	WaitAudioTrackDataMs    string `gorm:"column:wait_audio_track_data_ms" json:"wait_audio_track_data_ms"`
	WaitTrackReadyMs        string `gorm:"column:wait_track_ready_ms" json:"wait_track_ready_ms"`
	BroadcastPlayerCount    string `gorm:"column:broadcast_player_count_changed" json:"broadcast_player_count_changed"`
	CheckNvidiaDev          string `gorm:"column:check_nvidia_dev" json:"check_nvidia_dev"`
	MergeWriteMS            string `gorm:"column:merge_write_ms" json:"mergeWriteMS"`
	ListenIP                string `gorm:"column:listen_ip" json:"listen_ip"`
}

type HlsConfig struct {
	BroadcastRecordTs string `gorm:"column:hls_broadcast_record_ts" json:"broadcastRecordTs"`
	DeleteDelaySec    string `gorm:"column:hls_delete_delay_sec" json:"deleteDelaySec"`
	FastRegister      string `gorm:"column:hls_fast_register" json:"fastRegister"`
	FileBufSize       string `gorm:"column:hls_file_buf_size" json:"fileBufSize"`
	SegDelay          string `gorm:"column:hls_seg_delay" json:"segDelay"`
	SegDur            string `gorm:"column:hls_seg_dur" json:"segDur"`
	SegKeep           string `gorm:"column:hls_seg_keep" json:"segKeep"`
	SegNum            string `gorm:"column:hls_seg_num" json:"segNum"`
	SegRetain         string `gorm:"column:hls_seg_retain" json:"segRetain"`
}

type HookConfig struct {
	Alivestringerval     string `gorm:"column:alive_stringerval" json:"hook.alive_stringerval"`
	Enable               string `gorm:"column:enable" json:"hook.enable"`
	OnFlowReport         string `gorm:"column:on_flow_report" json:"hook.on_flow_report"`
	OnHttpAccess         string `gorm:"column:on_http_access" json:"hook.on_http_access"`
	OnPlay               string `gorm:"column:on_play" json:"hook.on_play"`
	OnPublish            string `gorm:"column:on_publish" json:"hook.on_publish"`
	OnRecordMp4          string `gorm:"column:on_record_mp4" json:"hook.on_record_mp4"`
	OnRecordTs           string `gorm:"column:on_record_ts" json:"hook.on_record_ts"`
	OnRtpServerTimeout   string `gorm:"column:on_rtp_server_timeout" json:"hook.on_rtp_server_timeout"`
	OnRtspAuth           string `gorm:"column:on_rtsp_auth" json:"hook.on_rtsp_auth"`
	OnRtspRealm          string `gorm:"column:on_rtsp_realm" json:"hook.on_rtsp_realm"`
	OnSendRtpStopped     string `gorm:"column:on_send_rtp_stopped" json:"hook.on_send_rtp_stopped"`
	OnServerExited       string `gorm:"column:on_server_exited" json:"hook.on_server_exited"`
	OnServerKeepalive    string `gorm:"column:on_server_keepalive" json:"hook.on_server_keepalive"`
	OnServerStarted      string `gorm:"column:on_server_started" json:"hook.on_server_started"`
	OnShellLogin         string `gorm:"column:on_shell_login" json:"hook.on_shell_login"`
	OnStreamChanged      string `gorm:"column:on_stream_changed" json:"hook.on_stream_changed"`
	OnStreamNoneReader   string `gorm:"column:on_stream_none_reader" json:"hook.on_stream_none_reader"`
	OnStreamNotFound     string `gorm:"column:on_stream_not_found" json:"hook.on_stream_not_found"`
	Retry                string `gorm:"column:retry" json:"hook.retry"`
	RetryDelay           string `gorm:"column:retry_delay" json:"hook.retry_delay"`
	StreamChangedSchemas string `gorm:"column:stream_changed_schemas" json:"hook.stream_changed_schemas"`
	TimeoutSec           string `gorm:"column:timeout_sec" json:"hook.timeoutSec"`
}

type HTTPConfig struct {
	AllowCrossDomains string `gorm:"column:allow_cross_domains" json:"http.allow_cross_domains"`
	AllowIPRange      string `gorm:"column:allow_ip_range" json:"http.allow_ip_range"`
	CharSet           string `gorm:"column:char_set" json:"http.charSet"`
	DirMenu           string `gorm:"column:dir_menu" json:"http.dirMenu"`
	ForbidCacheSuffix string `gorm:"column:forbid_cache_suffix" json:"http.forbidCacheSuffix"`
	ForwardedIPHeader string `gorm:"column:forwarded_ip_header" json:"http.forwarded_ip_header"`
	KeepAliveSecond   string `gorm:"column:keep_alive_second" json:"http.keepAliveSecond"`
	MaxReqSize        string `gorm:"column:max_req_size" json:"http.maxReqSize"`
	NotFound          string `gorm:"column:not_found" json:"http.notFound"`
	Port              string `gorm:"column:port" json:"http.port"`
	RootPath          string `gorm:"column:root_path" json:"http.rootPath"`
	SendBufSize       string `gorm:"column:send_buf_size" json:"http.sendBufSize"`
	SSLPort           string `gorm:"column:ssl_port" json:"http.sslport"`
	VirtualPath       string `gorm:"column:virtual_path" json:"http.virtualPath"`
}

type MulticastConfig struct {
	AddrMax string `gorm:"column:multicast_addr_max" json:"addrMax"`
	AddrMin string `gorm:"column:multicast_addr_min" json:"addrMin"`
	UdpTTL  string `gorm:"column:multicast_udp_ttl" json:"udpTTL"`
}

type ProtocolConfig struct {
	AddMuteAudio   string `json:"add_mute_audio" gorm:"column:add_mute_audio"`
	AutoClose      string `json:"auto_close" gorm:"column:auto_close"`
	ContinuePushMS string `json:"continue_push_ms" gorm:"column:continue_push_ms"`
	EnableAudio    string `json:"enable_audio" gorm:"column:enable_audio"`
	EnableFMP4     string `json:"enable_fmp4" gorm:"column:enable_fmp4"`
	EnableHLS      string `json:"enable_hls" gorm:"column:enable_hls"`
	EnableHLSFMP4  string `json:"enable_hls_fmp4" gorm:"column:enable_hls_fmp4"`
	EnableMP4      string `json:"enable_mp4" gorm:"column:enable_mp4"`
	EnableRTMP     string `json:"enable_rtmp" gorm:"column:enable_rtmp"`
	EnableRTSP     string `json:"enable_rtsp" gorm:"column:enable_rtsp"`
	EnableTS       string `json:"enable_ts" gorm:"column:enable_ts"`
	Fmp4Demand     string `json:"fmp4_demand" gorm:"column:fmp4_demand"`
	HlsDemand      string `json:"hls_demand" gorm:"column:hls_demand"`
	HlsSavePath    string `json:"hls_save_path" gorm:"column:hls_save_path"`
	ModifyStamp    string `json:"modify_stamp" gorm:"column:modify_stamp"`
	Mp4AsPlayer    string `json:"mp4_as_player" gorm:"column:mp4_as_player"`
	Mp4MaxSecond   string `json:"mp4_max_second" gorm:"column:mp4_max_second"`
	Mp4SavePath    string `json:"mp4_save_path" gorm:"column:mp4_save_path"`
	PacedSenderMS  string `json:"paced_sender_ms" gorm:"column:paced_sender_ms"`
	RtmpDemand     string `json:"rtmp_demand" gorm:"column:rtmp_demand"`
	RtspDemand     string `json:"rtsp_demand" gorm:"column:rtsp_demand"`
	TsDemand       string `json:"ts_demand" gorm:"column:ts_demand"`
}

type RecordConfig struct {
	AppName     string `gorm:"column:record_app_name" json:"appName"`
	EnableFmp4  string `gorm:"column:record_enable_fmp4" json:"enableFmp4"`
	FastStart   string `gorm:"column:record_fast_start" json:"fastStart"`
	FileBufSize string `gorm:"column:record_file_buf_size" json:"fileBufSize"`
	FileRepeat  string `gorm:"column:record_file_repeat" json:"fileRepeat"`
	SampleMS    string `gorm:"column:record_sample_ms" json:"sampleMS"`
}

type RTCConfig struct {
	BFilter              string `json:"bfilter" gorm:"column:rtc_bfilter"`
	DataChannelEcho      string `json:"datachannel_echo" gorm:"column:rtc_datachannel_echo"`
	ExternIP             string `json:"externIP" gorm:"column:rtc_extern_ip"`
	MaxRtpCacheMS        string `json:"maxRtpCacheMS" gorm:"column:rtc_max_rtp_cache_ms"`
	MaxRtpCacheSize      string `json:"maxRtpCacheSize" gorm:"column:rtc_max_rtp_cache_size"`
	MaxBitrate           string `json:"max_bitrate" gorm:"column:rtc_max_bitrate"`
	MinBitrate           string `json:"min_bitrate" gorm:"column:rtc_min_bitrate"`
	NackstringervalRatio string `json:"nackstringervalRatio" gorm:"column:rtc_nack_stringerval_ratio"`
	NackMaxCount         string `json:"nackMaxCount" gorm:"column:rtc_nack_max_count"`
	NackMaxMS            string `json:"nackMaxMS" gorm:"column:rtc_nack_max_ms"`
	NackMaxSize          string `json:"nackMaxSize" gorm:"column:rtc_nack_max_size"`
	NackRtpSize          string `json:"nackRtpSize" gorm:"column:rtc_nack_rtp_size"`
	Port                 string `json:"port" gorm:"column:rtc_port"`
	PreferredCodecA      string `json:"preferredCodecA" gorm:"column:rtc_preferred_codec_a"`
	PreferredCodecV      string `json:"preferredCodecV" gorm:"column:rtc_preferred_codec_v"`
	RembBitRate          string `json:"rembBitRate" gorm:"column:rtc_remb_bitrate"`
	StartBitrate         string `json:"start_bitrate" gorm:"column:rtc_start_bitrate"`
	TcpPort              string `json:"tcpPort" gorm:"column:rtc_tcp_port"`
	TimeoutSec           string `json:"timeoutSec" gorm:"column:rtc_timeout_sec"`
}

type RTMPConfig struct {
	DirectProxy     string `gorm:"column:rtmp_direct_proxy" json:"directProxy"`
	Enhanced        string `gorm:"column:rtmp_enhanced" json:"enhanced"`
	HandshakeSecond string `gorm:"column:rtmp_handshake_second" json:"handshakeSecond"`
	KeepAliveSecond string `gorm:"column:rtmp_keep_alive_second" json:"keepAliveSecond"`
	Port            string `gorm:"column:rtmp_port" json:"port"`
	SSLPort         string `gorm:"column:rtmp_ssl_port" json:"sslport"`
}

type RTPConfig struct {
	AudioMtuSize           string `json:"audioMtuSize" gorm:"column:rtp_audio_mtu_size"`
	H264StapA              string `json:"h264_stap_a" gorm:"column:rtp_h264_stap_a"`
	LowLatency             string `json:"lowLatency" gorm:"column:rtp_low_latency"`
	RtpMaxSize             string `json:"rtpMaxSize" gorm:"column:rtp_rtp_max_size"`
	VideoMtuSize           string `json:"videoMtuSize" gorm:"column:rtp_video_mtu_size"`
	RtpProxyDumpDir        string `json:"dumpDir" gorm:"column:rtp_proxy_dump_dir"`
	RtpProxyGopCache       string `json:"gop_cache" gorm:"column:rtp_proxy_gop_cache"`
	RtpProxyH264PT         string `json:"h264_pt" gorm:"column:rtp_proxy_h264_pt"`
	RtpProxyH265PT         string `json:"h265_pt" gorm:"column:rtp_proxy_h265_pt"`
	RtpProxyOpusPT         string `json:"opus_pt" gorm:"column:rtp_proxy_opus_pt"`
	RtpProxyPort           string `json:"port" gorm:"column:rtp_proxy_port"`
	RtpProxyPortRange      string `json:"port_range" gorm:"column:rtp_proxy_port_range"`
	RtpProxyPsPT           string `json:"ps_pt" gorm:"column:rtp_proxy_ps_pt"`
	RtpProxyRtpG711DurMs   string `json:"rtp_g711_dur_ms" gorm:"column:rtp_proxy_rtp_g711_dur_ms"`
	RtpProxyTimeoutSec     string `json:"timeoutSec" gorm:"column:rtp_proxy_timeout_sec"`
	RtpProxyUdpRecvSockBuf string `json:"udp_recv_socket_buffer" gorm:"column:rtp_proxy_udp_recv_socket_buffer"`
}

type RTSPConfig struct {
	AuthBasic        string `gorm:"column:rtsp_auth_basic" json:"authBasic"`
	DirectProxy      string `gorm:"column:rtsp_direct_proxy" json:"directProxy"`
	HandshakeSecond  string `gorm:"column:rtsp_handshake_second" json:"handshakeSecond"`
	KeepAliveSecond  string `gorm:"column:rtsp_keep_alive_second" json:"keepAliveSecond"`
	LowLatency       string `gorm:"column:rtsp_low_latency" json:"lowLatency"`
	Port             string `gorm:"column:rtsp_port" json:"port"`
	RtpTransportType string `gorm:"column:rtsp_rtp_transport_type" json:"rtpTransportType"`
	SSLPort          string `gorm:"column:rtsp_ssl_port" json:"sslport"`
}

type ShellConfig struct {
	MaxReqSize string `gorm:"column:shell_max_req_size" json:"maxReqSize"`
	Port       string `gorm:"column:shell_port" json:"port"`
}

type SrtConfig struct {
	LatencyMul string `json:"latencyMul" gorm:"column:srt_latency_mul"`
	PassPhrase string `json:"passPhrase" gorm:"column:srt_pass_phrase"`
	PktBufSize string `json:"pktBufSize" gorm:"column:srt_pkt_buf_size"`
	Port       string `json:"port" gorm:"column:srt_port"`
	TimeoutSec string `json:"timeoutSec" gorm:"column:srt_timeout_sec"`
}
