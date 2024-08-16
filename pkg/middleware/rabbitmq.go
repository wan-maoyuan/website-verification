package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"website-verification/pkg/conf"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var mq = Rabbit{}

type Rabbit struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	messageChan chan *Message
}

func InitRabbitmq() (err error) {
	config := conf.Get()

	mq.conn, err = amqp.Dial(config.MQURI)
	if err != nil {
		return fmt.Errorf("连接 rabbitmq 服务器失败: %v", err)
	}

	mq.channel, err = mq.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建一个 rabbitmq 连接通道失败: %v", err)
	}

	mq.messageChan = make(chan *Message)

	return nil
}

type Message struct {
	Url string `json:"url"`
	UA  string `json:"ua"`
}

func GetTaskMessage(ctx context.Context) (<-chan *Message, error) {
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
				var info = new(Message)

				if err := json.Unmarshal(message.Body, &info); err != nil {
					logrus.Errorf("rabbit mq 消息: %s 解析URL失败: %v", string(message.Body), err)
					continue
				}

				mq.messageChan <- info
				logrus.Debugf("从 rabbitmq 队列中收到一个信息: %s", info.Url)
			}
		}
	}()

	return mq.messageChan, nil
}

func SendMessage(msgChan chan *Message) error {
	publishQueue, err := mq.channel.QueueDeclare(
		conf.Get().MQQueue, // 队列名
		true,               // 是否持续
		false,              // 是否自动删除
		false,              // 是否独占
		false,              // 是否阻塞
		nil,                // args
	)

	if err != nil {
		return fmt.Errorf("获取 rabbit_mq 发布通道失败: %v", err)
	}

	for msg := range msgChan {
		body, err := json.Marshal(msg)
		if err != nil {
			logrus.Errorf("解析 Message 结构体 to json 失败: %v", err)
			continue
		}

		err = mq.channel.Publish(
			"",                // 交换机的名称
			publishQueue.Name, // 需要发送的消息队列
			false,             // 消息发送失败是否需要收到回复
			false,             // 设置为true，当消息无法直接投递到消费者时，会返回一个Basic.Return消息给生产者。如果设置为false，则消息会被存储在队列中，等待消费者连接。
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})

		if err != nil {
			logrus.Errorf("往 rabbit_mq 发送消息: %s 失败: %v", string(body), err)
			continue
		}
	}

	return nil
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
