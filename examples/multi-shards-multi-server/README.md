# multi-shards-multi-server

## 演示功能

在多个服务上，基于 redis session manager 启动多分片 websocket 连接。每个分片接收到的事件，将按照 guildID 进行 hash 分配。

相关文档：https://bot.q.qq.com/wiki/develop/api/gateway/shard.html

server1 和 server2 使用的是相同的`集群key`与`分配数`，基于 redis 会自动协同，进行分片的自动分配，处理不同的分片。
同时每次触发 resume 的时候，都会重新计算负责的分片。

当在相同的`集群key`上再启动一个进程的时候（总数 3 个），也会进入竞争，但是由于分片只有 2 个，所以会有一个进程闲置。
所以这里建议，分片数，要大于启动的进程数。就可以留有扩容的空间了。

## Redis Session Manager 代码

[remote](../../sessions/remote)