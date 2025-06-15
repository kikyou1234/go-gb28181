package m

import (
	"strings"
	"time"

	db "go-sip/db/sqlite"

	"github.com/spf13/viper"
)

// Config Config
type C_Config struct {
	UDP      string   `json:"udp" yaml:"udp" mapstructure:"udp"`
	Gateway  string   `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	LogLevel string   `json:"logger" yaml:"logger" mapstructure:"logger"`
	GB28181  *SysInfo `json:"gb28181" yaml:"gb28181" mapstructure:"gb28181"`
}

type SysInfo struct {
	// Region 当前域
	Region string `json:"region"   yaml:"region" mapstructure:"region"`
	// LID 当前服务id
	LID string `json:"lid" bson:"lid" yaml:"lid" mapstructure:"lid"`
	// 密码
	Passwd string `json:"passwd" bson:"passwd" yaml:"passwd" mapstructure:"passwd"`
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

}
