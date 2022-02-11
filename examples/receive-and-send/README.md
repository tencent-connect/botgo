# receive-and-send

## 演示功能

接收一条消息，根据消息的内容，识别指令，返回不同的消息内容。

### 指令类型

`time xxx`

返回文字信息，有关消息的发送时间，接受时间

`ark xxx`

返回一个 ark 消息

`dm xxx`

进入私信流程，机器人会给你发一条私信

## 启动方法

- 进入到 receive-and-send 目录
- 修改 config.yaml.demo 为 config.yaml，并在其中配置你的机器人 appid 和 token，注意要符合 yaml 格式要求
- 执行 `go build` 之后执行 `receive-and-send` 或者直接运行 `go run *.go`