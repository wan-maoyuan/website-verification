package main

import (
	"website-verification/pkg/conf"
	"website-verification/pkg/middleware"

	"github.com/sirupsen/logrus"
)

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	if err := middleware.InitRabbitmq(); err != nil {
		logrus.Errorf("InitRabbitmq: %v", err)
		return
	}
	defer middleware.CloseRabbitmq()

	ch := make(chan middleware.Message)

	go func() {
		defer close(ch)

		for i := 0; i < 10000; i++ {
			ch <- middleware.Message{Url: "https://www.baidu.com"}
			logrus.Info("发送了一条消息")
		}
	}()

	if err := middleware.SendMessage(ch); err != nil {
		logrus.Errorf("SendMessage: %v", err)
	}
}
