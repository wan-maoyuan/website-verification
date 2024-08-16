# website-verification
> 网页状态检查器，get 请求记录300和300以上的请求

## 代码编译打包流程
- 需要先下载 `golang` 版本不能低于 1.22。下载链接[https://golang.google.cn/dl/]

- `linux` 环境代码编译
```shell
git clone https://github.com/wan-maoyuan/website-verification.git
cd website-verification
make
```

## 代码运行
- 激活环境变量可以执行 
```shell
source resources/.env_dev.sh
```

- 环境变量解释
```shell
# rabbit_mq 链接 uri，必填
export MQ_URI="amqp://rabbit:123456@127.0.0.1:5672/"
# 需要监听的 rabbit_mq 队列名称，必填
export MQ_QUEUE="test"
# 请求处理的并发数，可以根据自己的电脑配置更改，默认是100个
export CONCURRENT=100
# 每一个请求的最大超时时间，默认10秒
export TIME_OUT_SECOND=10
# 网页处理计数器，统计的时间间隔，单位是分钟
export STATISTICS_INTERVAL_MINUTE=5
# 是否需要保存日志，为空不保存，直接在控制台打印。可以填 "./logs/website-verification.log"
export LOG_FILE=""
# 日志等级：debug,info,warn,error,fatal,panic
export LOG_LEVEL="info"
# 每一个日志文件最大的 size ，不超过 100 M
export LOG_SIZE=100
# 日志文件最大的个数，不超过100个
export LOG_AGE=100
```

- 运行程序
```shell
./dist/website-verification
```

## `docker` 镜像打包
```shell
make container
```

## `docker` 镜像通过 `docker-compose.yml` 文件运行
- 环境变量根据上面进行配置
```yml
version: '3.1'
services:
  rabbitmq:
    image: website-verification:v0.0.1
    container_name: website-verification
    environment:
      MQ_URI: "amqp://rabbit:123456@127.0.0.1:5672/"
      MQ_QUEUE: "test"
      CONCURRENT: 100
      TIME_OUT_SECOND: 10
      STATISTICS_INTERVAL_MINUTE: 5
      LOG_FILE: ""
      LOG_LEVEL: "info"
      LOG_SIZE: 100
      LOG_AGE: 100
```

- 启动 `docker` 镜像
```shell
docker-compose -f resources/docker-compose.yml up -d
```

- 停止 `docker` 容器
```shell
docker-compose -f resources/docker-compose.yml down
```

## 代码测试
- 在本地启动一个 `rabbit mq` 容器
```shell
# 启动
docker-compose -f resources/rabbitmq.yml up -d
# 关闭
docker-compose -f resources/rabbitmq.yml down
```

- 执行消息发送代码
> 网消息队列发送 10000 条消息
```shell
source resources/.env_dev.sh
go run test/main.go
```

- 启动主程序
```shell
make
source resources/.env_dev.sh
./dist/website-verification
```