# botgo
a golang sdk for guild bot

## 设计模式
分为三个主要模块

- openapi 用于请求 http 的 openapi
- websocket 用于监听事件网关，接收事件消息
- oauth 用于处理 oauth 的 token 获取

openapi 接口定义：`openapi/iface.go`，同时 sdk 中提供了 v1 的实现，后续 openapi 有新版本的时候，可以增加对应新版本的实现。
websocket 接口定义：`websocket/iface.go`，sdk 实现了默认版本的 client，如果开发者有更好的实现，也可以进行替换

## 使用

### 1.请求 openapi

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

### 2.请求 websocket

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
    var message websocket.MessageEventHandler = func(event *dto.WSMsg, data *dto.MessageData) error {
        fmt.Println(event, data)
        return nil
    }
    intent := websocket.RegisterHandlers(message)
    // 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
    botgo.Session.Start(ws, token, &intent)
}
```

### 3.请求 oauth 

待补充

### 4. session manager
接口定义：`session_manager.go`

sdk 实现了 `localSession` 主要是在单机上启动多个 shard 的连接，在实际生产中，如果需要启动多个 shard，那么有可能会采用分布式的管理方法，那么
就需要开发者自己实现一个分布式的 session manager 来进行连接管理。

## 开发说明

### 1. 如何增加新的 openapi 接口调用方法

- Step1: dto 中增加对应的对象
- Step2: openapi 的接口定义中，增加新方法的定义
- Step3：在 openapi 的实现中，实现这个新的方法

### 2. 如何增加新的 websocket 事件

- Step1: dto 中增加对应的对象 `dto/websocket_msg.go`
- Step2: 新增事件类型 `dto/websocket_msg.go`
- Step3: 新增 intent，以及事件对应的 intent（如果有）`dto/intents.go`
- Step4: 新增 event handler 类型，`websocket/event_handler.go`
- Step5：websocket 的具体实现中，针对收到的 message 进行解析，判断 type 是否符合新添加的时间类型，解析为 dto 之后，调用对应的 handler `websocket/client/event.go`

