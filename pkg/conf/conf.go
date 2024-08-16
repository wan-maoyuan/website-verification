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
	MQURI         string `mapstructure:"MQ_URI"`
	MQQueue       string `mapstructure:"MQ_QUEUE"`
	Concurrent    uint   `mapstructure:"CONCURRENT"`
	TimeoutSecond uint   `mapstructure:"TIME_OUT_SECOND"`
	Log           Log    `mapstructure:"LOG"`
}

func New() *Conf {
	viper.AutomaticEnv()

	config.MQURI = viper.GetString("MQ_URI")
	config.MQQueue = viper.GetString("MQ_QUEUE")
	config.Concurrent = viper.GetUint("CONCURRENT")
	config.TimeoutSecond = viper.GetUint("TIME_OUT_SECOND")

	config.Log = NewLog()

	config.Log.File = viper.GetString("LOG_FILE")
	config.Log.Level = viper.GetString("LOG_LEVEL")
	config.Log.MaxSize = viper.GetInt("LOG_SIZE")
	config.Log.MaxAge = viper.GetInt("LOG_AGE")

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
