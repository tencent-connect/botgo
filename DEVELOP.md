# 开发说明

## 一、SDK 的设计模式

分为三个主要模块

- openapi 用于请求 http 的 openapi
- websocket 用于监听事件网关，接收事件消息
- sessions 实现 session_manager 接口，用于管理 websocket 实例的新建，重连等

openapi 接口定义：`openapi/iface.go`，同时 sdk 中提供了 v1 的实现，后续 openapi 有新版本的时候，可以增加对应新版本的实现。
websocket 接口定义：`websocket/iface.go`，sdk 实现了默认版本的 client，如果开发者有更好的实现，也可以进行替换


## 二、SDK 增加新接口or新事件开发说明

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
