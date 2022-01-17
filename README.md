# BotGo

QQ频道机器人，官方 GOLANG SDK。

![Build](https://github.com/tencent-connect/botgo/actions/workflows/build.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/tencent-connect/botgo.svg)](https://pkg.go.dev/github.com/tencent-connect/botgo)
[![Examples](https://img.shields.io/badge/BotGo-examples-yellowgreen)](https://github.com/tencent-connect/botgo/tree/master/examples)


## 一、如何使用

### 1.请求 openapi 接口，操作资源

```golang
func main() {
    token := token.BotToken(conf.AppID, conf.Token)
    api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
    ctx := context.Background()

    ws, err := api.WS(ctx, nil, "")
    log.Printf("%+v, err:%v", ws, err)
    
    me, err := api.Me(ctx, nil, "")
    log.Printf("%+v, err:%v", me, err)
}
```

### 2.使用默认 SessionManager 启动 websocket 连接，接收事件

```golang
func main() {
    token := token.BotToken(conf.AppID, conf.Token)
    api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
    ctx := context.Background()
    ws, err := api.WS(ctx, nil, "")
    if err != nil {
        log.Printf("%+v, err:%v", ws, err)
    }

    // 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
    var atMessage websocket.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
        fmt.Println(event, data)
        return nil
    }
    intent := websocket.RegisterHandlers(atMessage)
    // 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
    botgo.NewSessionManager().Start(ws, token, &intent)
}
```

## 二、什么是 SessionManager

SessionManager，用于管理 websocket 连接的启动，重连等。接口定义在：`session_manager.go`。开发者也可以自己实现自己的 SessionManager。

sdk 中实现了两个 SessionManager

- [local](./sessions/local/local.go) 用于在单机上启动多个 shard 的连接。下文用 `local` 代表
- [remote](./sessions/remote/remote.go) 基于 redis 的 list 数据结构，实现分布式的 shard 管理，可以在多个节点上启动多个服务进程。下文用 `remote` 代表

另外，也有其他同事基于 etcd 实现了 shard 集群的管理，在 [botgo-plugns](https://github.com/tencent-connect/botgo-plugins) 中。

## 三、生产环境中的一些建议

得益于 websocket 的机制，我们可以在本地就启动一个机器人，实现相关逻辑，但是在生产环境中需要考虑扩容，容灾等情况，所以建
议从以下几方面考虑生产环境的部署：

### 1.公域机器人，优先使用分布式 shard 管理

使用上面提到的分布式的 session manager 或者自己实现一个分布式的 session manager

### 2.提前规划好分片

分布式 SessionManager 需要解决的最大的问题，就是如何解决 shard 随时增加的问题，类似 kafka 的 rebalance 问题一样，
由于 shard 是基于频道 id 来进行 hash 的，所以在扩容的时候所有的数据都会被重新 hash。

提前规划好较多的分片，如 20 个分片，有助于在未来机器人接入的频道过多的时候，能够更加平滑的进行实例的扩容。比如如果使用的
是 `remote`，初始化时候分 20 个分片，但是只启动 2 个进程，那么这2个进程将争抢 20 个分片的消费权，进行消费，当启动更多
的实例之后，伴随着 websocket 要求一定时间进行一次重连，启动的新实例将会平滑的分担分片的数据处理。

### 3.接入和逻辑分离

接入是指从机器人平台收到事件的服务。逻辑是指处理相关事件的服务。

接入与逻辑分离，有助于提升机器人的事件处理效率和可靠性。一般实现方式类似于以下方案：

- 接入层：负责维护与平台的 websocket 连接，并接收相关事件，生产到 kafka 等消息中间件中。
  如果使用 `local` 那么可能还涉及到分布式锁的问题。可以使用sdk 中的 `sessions/remote/lock` 快速基于 redis 实现分布式锁。

- 逻辑层：从 kafka 消费到事件，并进行对应的处理，或者调用机器人的 openapi 进行相关数据的操作。

提前规划好 kafka 的分片，然后从容的针对逻辑层做水平扩容。或者使用 pulsar（腾讯云上叫 tdmq） 来替代 kafka 避免 rebalance 问题。

## 四、SDK 开发说明

请查看：[开发说明](./DEVELOP.md)

## 五、加入官方社区

欢迎扫码加入 **QQ 频道开发者社区**。

![开发者社区](https://mpqq.gtimg.cn/privacy/qq_guild_developer.png)
