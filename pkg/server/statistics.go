package server

import (
	"context"
	"sync/atomic"
	"time"
	"website-verification/pkg/conf"

	"github.com/sirupsen/logrus"
)

type AtomicCounter struct {
	successCount int64
	failCount    int64
	duration     time.Duration
	ticker       *time.Ticker
}

func NewAtomicCounter() *AtomicCounter {
	var duration time.Duration
	if conf.Get().StatisticsIntervalMinute < 1 {
		duration = time.Minute
	} else {
		duration = time.Minute * time.Duration(conf.Get().StatisticsIntervalMinute)
	}

	return &AtomicCounter{
		successCount: 0,
		failCount:    0,
		ticker:       time.NewTicker(duration),
		duration:     duration,
	}
}

func (counter *AtomicCounter) Run(ctx context.Context) {
	logrus.Infof("计时器开始工作")

	for {
		select {
		case <-ctx.Done():
			logrus.Infof("计时器停止")
			return
		case <-counter.ticker.C:
			logrus.Infof("%v 时间内网页检测成功的个数: %d", counter.duration, counter.successCount)
			logrus.Infof("%v 时间内网页检测失败的个数: %d", counter.duration, counter.failCount)
		}
	}
}

func (counter *AtomicCounter) AddSucces() {
	atomic.AddInt64(&counter.successCount, 1)
}

func (counter *AtomicCounter) AddFail() {
	atomic.AddInt64(&counter.failCount, 1)
}
