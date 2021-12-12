# multi-shards-single-server

## 演示功能

单服务启动多分片 websocket 连接。每个分片接收到的事件，将按照 guildID 进行 hash 分配。

相关文档：https://bot.q.qq.com/wiki/develop/api/gateway/shard.html