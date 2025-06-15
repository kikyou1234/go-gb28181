package m

import (
	"strings"
	"time"

	db "go-sip/db/sqlite"

	"fmt"

	"github.com/spf13/viper"
)

// Config Config
type C_Config struct {
	API           string       `json:"api" yaml:"api" mapstructure:"api"`
	SipClientPort string       `json:"sip_client_port" yaml:"sip_client_port" mapstructure:"sip_client_port"`
	UDP           string       `json:"udp" yaml:"udp" mapstructure:"udp"`
	TCP           string       `json:"tcp" yaml:"tcp" mapstructure:"tcp"`
	Gateway       string       `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	ZlmSecret     string       `json:"zlm_secret" yaml:"zlm_secret" mapstructure:"zlm_secret"`
	ZlmInnerIp    string       `json:"zlm_inner_ip" yaml:"zlm_inner_ip" mapstructure:"zlm_inner_ip"`
	LogLevel      string       `json:"logger" yaml:"logger" mapstructure:"logger"`
	Stream        *Stream      `json:"stream" yaml:"stream" mapstructure:"stream"`
	GB28181       *SysInfo     `json:"gb28181" yaml:"gb28181" mapstructure:"gb28181"`
	Audio         *AudioConfig `json:"audio" yaml:"audio" mapstructure:"audio"`
}

// Stream Stream
type Stream struct {
	HLS  bool `json:"hls" yaml:"hls" mapstructure:"hls"`
	RTMP bool `json:"rtmp" yaml:"rtmp" mapstructure:"rtmp"`
}

type SysInfo struct {
	// Region 当前域
	Region string `json:"region"   yaml:"region" mapstructure:"region"`
	// LID 当前服务id
	LID string `json:"lid" bson:"lid" yaml:"lid" mapstructure:"lid"`
	// 密码
	Passwd string `json:"passwd" bson:"passwd" yaml:"passwd" mapstructure:"passwd"`
}

type AudioConfig struct {
	SampleRate   float64 `json:"sample_rate" yaml:"sample_rate" mapstructure:"sample_rate"`
	Channels     int     `json:"channels" yaml:"channels" mapstructure:"channels"`
	InputDevice  int     `json:"input_device" yaml:"input_device" mapstructure:"input_device"`
	OutputDevice int     `json:"output_device" yaml:"output_device" mapstructure:"output_device"`
}

func DefaultInfo() *SysInfo {
	return CMConfig.GB28181
}

var CMConfig *C_Config

func LoadClientConfig() {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetDefault("logger", "debug")
	viper.SetDefault("mod", "release")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	CMConfig = &C_Config{}
	err = viper.Unmarshal(&CMConfig)
	if err != nil {
		panic(err)
	}
	db.DBClient, err = db.Open()
	if err != nil {
		panic(err)
	}
	db.DBClient.SetNowFuncOverride(func() interface{} {
		return time.Now().Unix()
	})
	db.DBClient.LogMode(false)
	go db.KeepLive(db.DBClient, time.Minute)

	// 自动设置 API 为 0.0.0.0:<sip_port>（如果未设置）
	api := strings.TrimSpace(CMConfig.API)
	if api == "" || api == "0.0.0.0" || api == "0.0.0.0:" {
		CMConfig.API = fmt.Sprintf("0.0.0.0:%s", CMConfig.SipClientPort)
	}

}
