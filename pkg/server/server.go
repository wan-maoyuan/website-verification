package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
	"website-verification/pkg/conf"
	"website-verification/pkg/middleware"

	"github.com/sirupsen/logrus"
)

type Verificationer struct {
	ch         chan struct{}
	httpClient *http.Client
	counter    *AtomicCounter
}

func NewVerificationer() (*Verificationer, error) {
	return &Verificationer{
		ch: make(chan struct{}, conf.Get().Concurrent),
		httpClient: &http.Client{
			Timeout: time.Duration(conf.Get().TimeoutSecond) * time.Second,
		},
		counter: NewAtomicCounter(),
	}, nil
}

func (srv *Verificationer) Run(ctx context.Context) error {
	// 启动计数器
	go srv.counter.Run(ctx)

	ch, err := middleware.GetTaskMessage(ctx)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

a:
	for {
		select {
		case <-ctx.Done():
			logrus.Info("程序停止运行")
			break a
		case message := <-ch:
			if message != nil {
				srv.handlerUrl(message, wg)
			}
		}
	}

	wg.Wait()
	return nil
}

func (srv *Verificationer) handlerUrl(message *middleware.Message, wg *sync.WaitGroup) {
	srv.ch <- struct{}{}
	wg.Add(1)

	go srv.verificationUrl(message, wg)
}

func (srv *Verificationer) verificationUrl(message *middleware.Message, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-srv.ch
	}()

	urlStr := fmt.Sprintf("%s&ua=%s", message.Url, url.QueryEscape(message.UA))
	resp, err := srv.httpClient.Get(urlStr)
	if err != nil {
		logrus.Errorf("请求处理失败: %v 网址: %s ua: %s", err, message.Url, message.UA)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusMultipleChoices {
		srv.counter.AddSucces()
		logrus.Debugf("请求处理成功，网址: %s ua: %s", message.Url, message.UA)
	} else {
		srv.counter.AddFail()
		logrus.Errorf("请求处理失败: %v 网址: %s ua: %s 网址状态响应码: %d", err, message.Url, message.UA, resp.StatusCode)
	}
}
