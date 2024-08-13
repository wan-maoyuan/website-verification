package server

import (
	"context"
	"fmt"
	"website-verification/pkg/middleware"

	"github.com/sirupsen/logrus"
)

type Verificationer struct {
}

func NewVerificationer() (*Verificationer, error) {

	return nil, nil
}

func (srv *Verificationer) Run(ctx context.Context) error {
	ch, err := middleware.GetTaskMessage(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Info("程序停止运行")
			return nil
		case url := <-ch:
			fmt.Println("url", url)
		}
	}
}
