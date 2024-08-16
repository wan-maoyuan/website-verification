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

	ch := make(chan *middleware.Message)

	go func() {
		defer close(ch)

		for i := 0; i < 100; i++ {
			ch <- &middleware.Message{
				Url: "http://47.239.194.150:31500/api/com.hello.jim?pid=jimobi&clickid=8647d5d5-27dd-4920-b799-e090670cd6e5&ip=216.122.165.214&android_id=95335533-c65b-474e-8a89-85ba722a5f13",
				UA:  "Mozilla/5.0 (Linux; Android 14; V2323A Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/116.0.0.0 Mobile Safari/537.36",
			}
		}
	}()

	if err := middleware.SendMessage(ch); err != nil {
		logrus.Errorf("SendMessage: %v", err)
	}
}
