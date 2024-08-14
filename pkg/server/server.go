package server

import (
	"context"
	"net/http"
	"sync"
	"time"
	"website-verification/pkg/conf"
	"website-verification/pkg/middleware"

	"github.com/sirupsen/logrus"
)

type Verificationer struct {
	ch         chan struct{}
	httpClient *http.Client
}

func NewVerificationer() (*Verificationer, error) {
	return &Verificationer{
		ch: make(chan struct{}, conf.Get().Concurrent),
		httpClient: &http.Client{
			Timeout: time.Duration(conf.Get().TimeoutSecond) * time.Second,
		},
	}, nil
}

func (srv *Verificationer) Run(ctx context.Context) error {
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
		case url := <-ch:
			srv.handlerUrl(url, wg)
		}
	}

	wg.Wait()
	return nil
}

func (srv *Verificationer) handlerUrl(url string, wg *sync.WaitGroup) {
	srv.ch <- struct{}{}
	wg.Add(1)

	go srv.verificationUrl(url, wg)
}

func (srv *Verificationer) verificationUrl(url string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-srv.ch
	}()

	resp, err := srv.httpClient.Get(url)
	if err != nil {
		logrus.Errorf("请求处理失败: %v 网址: %s", err, url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusMultipleChoices {
		logrus.Debugf("请求处理成功，网址: %s", url)
	} else {
		logrus.Errorf("请求处理失败: %v 网址: %s 网址状态响应码: %d", err, url, resp.StatusCode)
	}
}
