package m

import (
	"strings"

	"github.com/spf13/viper"
)

// Config Config
type G_Config struct {
	GatewayAPI string      `json:"gateway_api" yaml:"gateway_api" mapstructure:"gateway_api"`
	Sign       string      `json:"sign" yaml:"sign" mapstructure:"sign"`
	DateBase   RedisConfig `json:"database" yaml:"database" mapstructure:"database"`
}

var GatewayConfig *G_Config

func LoadGatewsyConfig() {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetDefault("gateway_api", "0.0.0.0:8999")
	viper.SetDefault("mod", "release")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	GatewayConfig = &G_Config{}
	err = viper.Unmarshal(&GatewayConfig)
	if err != nil {
		panic(err)
	}
}
