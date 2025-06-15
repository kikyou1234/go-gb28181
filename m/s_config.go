package m

import (
	"strings"

	"github.com/spf13/viper"
)

// Config Config
type S_Config struct {
	SipID      string      `json:"sip_id" yaml:"sip_id" mapstructure:"sip_id"`
	Api        ApiConfig   `json:"api" yaml:"api" mapstructure:"api"`
	Sip        SipConfig   `json:"sip" yaml:"sip" mapstructure:"sip"`
	GatewayAPI string      `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Secret     string      `json:"secret" yaml:"secret" mapstructure:"secret"`
	Sign       string      `json:"sign" yaml:"sign" mapstructure:"sign"`
	DateBase   RedisConfig `json:"database" yaml:"database" mapstructure:"database"`
}

type ApiConfig struct {
	IP   string `json:"ip" yaml:"ip" mapstructure:"ip"`
	Port string `json:"port" yaml:"port" mapstructure:"port"`
}

type SipConfig struct {
	SipIP   string `json:"ip" yaml:"ip" mapstructure:"ip"`
	SipPort string `json:"port" yaml:"port" mapstructure:"port"`
}

type RedisConfig struct {
	Dialect  string `json:"dialect" yaml:"dialect" mapstructure:"dialect"`
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Password string `json:"password" yaml:"passwintord" mapstructure:"password"`
	DB       int    `json:"DB" yaml:"DB" mapstructure:"DB"`
}

var SMConfig *S_Config

func LoadServerConfig() {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetDefault("mod", "release")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	SMConfig = &S_Config{}
	err = viper.Unmarshal(&SMConfig)
	if err != nil {
		panic(err)
	}

}
