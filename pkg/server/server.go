package server

import (
	"context"
	"net/http"
	"net/url"
	"strings"
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
		case urlStr := <-ch:
			srv.handlerUrl(urlStr, wg)
		}
	}

	wg.Wait()
	return nil
}

func (srv *Verificationer) handlerUrl(urlStr string, wg *sync.WaitGroup) {
	srv.ch <- struct{}{}
	wg.Add(1)

	go srv.verificationUrl(urlStr, wg)
}

func (srv *Verificationer) verificationUrl(urlStr string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-srv.ch
	}()

	resp, err := srv.httpClient.Get(encodingUrl(urlStr))
	if err != nil {
		logrus.Errorf("请求处理失败: %v 网址: %s", err, urlStr)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusMultipleChoices {
		logrus.Debugf("请求处理成功，网址: %s", urlStr)
	} else {
		logrus.Errorf("请求处理失败: %v 网址: %s 网址状态响应码: %d", err, urlStr, resp.StatusCode)
	}
}

// 对 URL 进行编码
func encodingUrl(urlStr string) string {
	// 查找第一个问号的位置
	index := strings.Index(urlStr, "?")
	if index == -1 {
		// 不包含参数，不需要进行 url 编码
		return urlStr
	}

	// 基础URL地址
	baseUrl := urlStr[:index] + "?"

	// 截取问号后面的部分
	paramStr := url.QueryEscape(urlStr[index+1:])

	return baseUrl + paramStr
}
