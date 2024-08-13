package middleware

import (
	"context"
	"fmt"
	"website-verification/pkg/conf"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var mq = Rabbit{}

type Rabbit struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	messageChan chan string
}

func InitRabbitmq() (err error) {
	config := conf.Get()

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.MQUser, config.MQPwd, config.MQHost, config.MQPort)
	mq.conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("连接 rabbitmq 服务器失败: %v", err)
	}

	mq.channel, err = mq.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建一个 rabbitmq 连接通道失败: %v", err)
	}

	mq.messageChan = make(chan string)

	return nil
}

func GetTaskMessage(ctx context.Context) (<-chan string, error) {
	deliveryChan, err := mq.channel.Consume(
		conf.Get().MQQueue, // 队列名
		"",                 // 消费者标签
		true,               // 是否自动回复
		false,              // 是否独占
		true,               // 是否阻塞
		false,              // 其他属性
		nil,                // args
	)

	if err != nil {
		return mq.messageChan, fmt.Errorf("rabbit mq 监听队列: %s 失败: %v", conf.Get().MQQueue, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				logrus.Infof("停止监听 rabbit mq 消息队列: %s", conf.Get().MQQueue)
				return
			case message := <-deliveryChan:
				url := string(message.Body)
				mq.messageChan <- url

				logrus.Debugf("从 rabbitmq 队列中收到一个信息: %s", url)
			}
		}
	}()

	return mq.messageChan, nil
}

func CloseRabbitmq() {
	if mq.messageChan != nil {
		close(mq.messageChan)
	}

	if mq.channel != nil {
		if err := mq.channel.Close(); err != nil {
			logrus.Errorf("关闭 rabbitmq 连接通道失败: %v", err)
		}
	}

	if mq.conn != nil {
		if err := mq.conn.Close(); err != nil {
			logrus.Errorf("关闭 rabbitmq 连接失败: %v", err)
		}
	}
}
