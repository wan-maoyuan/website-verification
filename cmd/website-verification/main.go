package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"website-verification/pkg/conf"
	"website-verification/pkg/middleware"
	"website-verification/pkg/server"

	"github.com/sirupsen/logrus"
)

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	if err := BeforeStartFunc(); err != nil {
		logrus.Errorf("服务初始化失败: %v", err)
		return
	}

	server, err := server.NewVerificationer()
	if err != nil {
		logrus.Errorf("服务创建失败: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())
	finishCh := make(chan os.Signal, 1)

	go func() {
		if err := server.Run(ctx); err != nil {
			logrus.Errorf("服务启动失败: %v", err)
			finishCh <- syscall.SIGQUIT
		}
	}()

	signal.Notify(finishCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-finishCh

	cancel()
	AfterStopFunc()
	logrus.Info("服务停止")
}

func BeforeStartFunc() error {
	if err := middleware.InitRabbitmq(); err != nil {
		return fmt.Errorf("初始化 rabbit mq 失败: %v", err)
	}

	return nil
}

func AfterStopFunc() error {
	middleware.CloseRabbitmq()

	return nil
}
