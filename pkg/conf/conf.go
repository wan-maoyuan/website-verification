package conf

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var config = new(Conf)

func Get() *Conf {
	return config
}

type Conf struct {
	MQHost     string `mapstructure:"MQ_HOST"`
	MQPort     string `mapstructure:"MQ_PORT"`
	MQUser     string `mapstructure:"MQ_USER"`
	MQPwd      string `mapstructure:"MQ_PWD"`
	MQQueue    string `mapstructure:"MQ_QUEUE"`
	Concurrent uint   `mapstructure:"CONCURRENT"`
	Log        Log    `mapstructure:"LOG"`
}

func New() *Conf {
	viper.AutomaticEnv()

	config.MQHost = viper.GetString("MQ_HOST")
	config.MQPort = viper.GetString("MQ_PORT")
	config.MQUser = viper.GetString("MQ_USER")
	config.MQPwd = viper.GetString("MQ_PWD")
	config.MQQueue = viper.GetString("MQ_QUEUE")
	config.Concurrent = viper.GetUint("CONCURRENT")

	config.Log.File = viper.GetString("LOG_FILE")
	config.Log.Level = viper.GetString("LOG_LEVEL")

	config.Log = NewLog()
	config.Log.InitLog()

	return config
}

func (c *Conf) Show() {
	if b, err := yaml.Marshal(c); err != nil {
		return
	} else {
		fmt.Printf(`
-----------------------------------------------------------------------------------------
%s
-----------------------------------------------------------------------------------------
`, string(b))
	}
}
