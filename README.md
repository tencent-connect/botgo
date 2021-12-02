# botgo

a golang sdk for guild bot

## 设计模式

分为三个主要模块

- openapi 用于请求 http 的 openapi
- websocket 用于监听事件网关，接收事件消息
- sessions 实现 session_manager 接口，用于管理 websocket 实例的新建，重连等

openapi 接口定义：`openapi/iface.go`，同时 sdk 中提供了 v1 的实现，后续 openapi 有新版本的时候，可以增加对应新版本的实现。
websocket 接口定义：`websocket/iface.go`，sdk 实现了默认版本的 client，如果开发者有更好的实现，也可以进行替换

## 使用

### 1.请求 openapi 接口

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

### 2.启动一个单独的 websocket 连接

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

### 3. 使用 session manager 启动多个连接

接口定义在：`session_manager.go`

sdk 实现了两个 session manager

- [local](./sessions/local/local.go) 用于在单机上启动多个 shard 的连接。下文用 `local` 代表
- [remote](./sessions/remote/remote.go) 基于 redis 的 list 数据结构，实现分布式的 shard 管理，可以在多个节点上启动多个服务进程。下文用 `remote` 代表

另外，也有其他同事基于 etcd 实现了 shard 集群的管理，在 [botgo-plugns](https://github.com/tencent-connect/botgo-plugins) 中。

### 4. 生产环境中的建议

得益于 websocket 的机制，我们可以在本地就启动一个机器人，实现相关逻辑，但是在生产环境中需要考虑扩容，容灾等情况，所以建议从以下几方面考虑生产环境的部署：

#### 公域机器人，优先使用分布式 shard 管理

使用上面提到的分布式的 session manager 或者自己实现一个分布式的 session manager

#### 提前规划好分片

分布式 session manager 需要解决的最大的问题，就是如何解决 shard 随时增加的问题，类似 kafka 的 rebalance 问题一样，由于 shard 是基于频道 id 来进行 hash 的，所以在扩容的时候所有的数据都会被重新 hash。

提前规划好较多的分片，如 20 个分片，有助于在未来机器人接入的频道过多的时候，能够更加平滑的进行实例的扩容。比如如果使用的是 `remote`，初始化时候分 20 个分片，但是只启动 2 个进程，那么这2个进程将争抢 20 个分片的消费权，进行消费，当启动更多的实例之后，伴随着 websocket 要求一定时间进行一次重连，启动的新实例将会平滑的分担分片的数据处理。

#### 接入和逻辑分离

接入是指从机器人平台收到事件的服务。逻辑是指处理相关事件的服务。

接入与逻辑分离，有助于提升机器人的事件处理效率和可靠性。一般实现方式类似于以下方案：

- 接入层：负责维护与平台的 websocket 连接，并接收相关事件，生产到 kafka 等消息中间件中。
  如果使用 `local` 那么可能还涉及到分布式锁的问题。可以使用sdk 中的 `sessions/remote/lock` 快速基于 redis 实现分布式锁。

- 逻辑层：从 kafka 消费到事件，并进行对应的处理，或者调用机器人的 openapi 进行相关数据的操作。

提前规划好 kafka 的分片，然后从容的针对逻辑层做水平扩容。或者使用 pulsar（腾讯云上叫 tdmq） 来替代 kafka 避免 rebalance 问题。

## SDK 增加新接口or新事件开发说明

### 1. 如何增加新的 openapi 接口调用方法（预计耗时3min）

- Step1: dto 中增加对应的对象
- Step2: openapi 的接口定义中，增加新方法的定义
- Step3：在 openapi 的实现中，实现这个新的方法

### 2. 如何增加新的 websocket 事件（预计耗时10min）

- Step1: dto 中增加对应的对象 `dto/websocket_payload.go`
- Step2: 新增 intent，以及事件对应的 intent（如果有）`dto/intents.go`
- Step3: 新增事件类型与 intent 的关系 `dto/websocket_event.go`
- Step4: 新增 event handler 类型，并在注册方法中补充断言，`websocket/event_handler.go`
- Step5：websocket 的具体实现中，针对收到的 message 进行解析，判断 type 是否符合新添加的时间类型，解析为 dto 之后，调用对应的 handler `websocket/client/event.go`
